package api

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/png"
	"mime/multipart"
	"net"
	"net/http"
	"net/textproto"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/healthy-heroes/neskuchka/backend/app/internal/testutil"
)

// The ceiling above the API has to stay above the longest route under it.
// Should it ever drop below, the upload would lose the budget it asks for
// without saying so — chi's Timeout silently keeps whichever deadline is
// earlier.
func Test_AvatarUploadFitsUnderTheCeiling(t *testing.T) {
	require.Less(t, avatarUploadTimeout, maxRequestTimeout)
}

// A slow upload used to be read in full, decoded, and only then thrown away:
// the route sat inside the budget every other route shares, so the write at the
// end met a context that had expired while the body was still arriving, and the
// client got a 500 for a picture the server had already processed.
//
// The route has a budget of its own now. Nothing cheaper notices if it loses
// one again — the constants can't, since the budget is derived from them, and
// the fault is in how the router is composed rather than in what it composes.
// Hence a real connection, and hence the seconds this test costs.
func Test_UploadAvatar_SurvivesASlowConnection(t *testing.T) {
	app := NewTestApp(t)

	user, err := app.DataStorage.CreateUser(t.Context(), testutil.CreateUser())
	require.NoError(t, err)
	cookie := app.LoginAs(t, user.ID)

	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	encoded := new(bytes.Buffer)
	require.NoError(t, png.Encode(encoded, img))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="avatar"; filename="photo.png"`)
	header.Set("Content-Type", "image/png")

	part, err := writer.CreatePart(header)
	require.NoError(t, err)
	_, err = part.Write(encoded.Bytes())
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	payload := body.Bytes()
	half := len(payload) / 2

	conn, err := net.Dial("tcp", app.Server.Listener.Addr().String())
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })

	_, err = fmt.Fprintf(conn, "POST /api/v1/user/me/avatar HTTP/1.1\r\nHost: test\r\n"+
		"Cookie: %s=%s\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n",
		cookie.Name, cookie.Value, writer.FormDataContentType(), len(payload))
	require.NoError(t, err)

	_, err = conn.Write(payload[:half])
	require.NoError(t, err)

	// Long enough that a route on the shared budget would already have given
	// up, short enough to stay well inside this one.
	time.Sleep(requestTimeout + time.Second)

	_, err = conn.Write(payload[half:])
	require.NoError(t, err)

	require.NoError(t, conn.SetReadDeadline(time.Now().Add(10*time.Second)))
	resp, err := http.ReadResponse(bufio.NewReader(conn), nil)
	require.NoError(t, err)
	t.Cleanup(func() { resp.Body.Close() })

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	stored, err := app.AvatarStorage.Get(t.Context(), user.ID)
	require.NoError(t, err, "the upload answered OK but stored nothing")
	assert.NotEmpty(t, stored.Data)
}
