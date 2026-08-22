package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// The throttle's own behaviour is covered in middlewares; what belongs here is
// the one thing that ties it to this router. Waiting longer than the route
// allows would hand the client chi's timeout instead of the 503 the throttle
// means to send, and the decode itself still has to fit in what is left.
func Test_AvatarUploadWaitFitsRouteTimeout(t *testing.T) {
	require.Less(t, avatarUploadWait, requestTimeout)
}
