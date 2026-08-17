package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/testutil"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/token"
	"github.com/healthy-heroes/neskuchka/backend/app/storage/datastorage"
)

const (
	testIssuer = "Neskuchka"
	testSecret = "test_secret"
)

type sentEmail struct {
	To      string
	Subject string
	Text    string
}

type fakeEmailSender struct {
	sent []sentEmail
	err  error
}

func (f *fakeEmailSender) Send(to, subject, text string) error {
	if f.err != nil {
		return f.err
	}

	f.sent = append(f.sent, sentEmail{To: to, Subject: subject, Text: text})
	return nil
}

// fakeTemplater records the tokens it was asked to template,
// so a test can follow the exact token that went out in the email
type fakeTemplater struct {
	tokens []string
	err    error
}

func (f *fakeTemplater) AuthLink(token string) (string, error) {
	if f.err != nil {
		return "", f.err
	}

	f.tokens = append(f.tokens, token)
	return "auth link for " + token, nil
}

type fakeSessionManager struct {
	setUserIDs []string
	cleared    int
	err        error
}

func (f *fakeSessionManager) Set(_ http.ResponseWriter, userID string) error {
	if f.err != nil {
		return f.err
	}

	f.setUserIDs = append(f.setUserIDs, userID)
	return nil
}

func (f *fakeSessionManager) Clear(_ http.ResponseWriter) {
	f.cleared++
}

type authFixture struct {
	Service   *Service
	Storage   *datastorage.Storage
	Email     *fakeEmailSender
	Templater *fakeTemplater
	Session   *fakeSessionManager
}

func setupAuth(t *testing.T) *authFixture {
	t.Helper()

	storage := datastorage.New(testutil.NewEngine(t), zerolog.Nop())

	f := &authFixture{
		Storage:   storage,
		Email:     &fakeEmailSender{},
		Templater: &fakeTemplater{},
		Session:   &fakeSessionManager{},
	}

	f.Service = NewService(
		domain.NewStore(domain.Opts{Storage: storage}),
		f.Session,
		Opts{
			Issuer: testIssuer,
			Secret: testSecret,
			Logger: zerolog.Nop(),

			EmailSender:    f.Email,
			EmailTemplater: f.Templater,
		},
	)

	return f
}

func (f *authFixture) call(t *testing.T, handler http.HandlerFunc, body any) *httptest.ResponseRecorder {
	t.Helper()

	buf, err := json.Marshal(body)
	require.NoError(t, err)

	r := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", bytes.NewReader(buf))
	r.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler(w, r)

	return w
}

// login runs the login handler and returns the token that reached the email
func (f *authFixture) login(t *testing.T, email string) string {
	t.Helper()

	w := f.call(t, f.Service.Login, LoginSchema{Email: email})
	require.Equal(t, http.StatusOK, w.Code)
	require.Len(t, f.Templater.tokens, 1)

	return f.Templater.tokens[0]
}

// makeToken mints a confirmation token outside the service, to forge the cases
// a well-behaved client never produces: a wrong secret, an expired deadline
func makeToken(t *testing.T, secret, email string, expiresAt time.Time, jti string) string {
	t.Helper()

	claims := jwt.MapClaims{
		"iss":  testIssuer,
		"jti":  jti,
		"exp":  expiresAt.Unix(),
		"data": map[string]string{"email": email},
	}

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	require.NoError(t, err)

	return signed
}

func Test_AuthService_Login(t *testing.T) {
	t.Run("should email a token carrying the address the user logs in with", func(t *testing.T) {
		f := setupAuth(t)

		w := f.call(t, f.Service.Login, LoginSchema{Email: "user@example.com"})
		assert.Equal(t, http.StatusOK, w.Code)

		require.Len(t, f.Email.sent, 1)
		assert.Equal(t, "user@example.com", f.Email.sent[0].To)
		assert.Equal(t, "Вход в Нескучку", f.Email.sent[0].Subject)

		require.Len(t, f.Templater.tokens, 1)
		assert.Equal(t, "auth link for "+f.Templater.tokens[0], f.Email.sent[0].Text)

		var claims struct {
			jwt.RegisteredClaims
			Data LoginSchema `json:"data"`
		}
		err := token.NewService(token.Opts{Issuer: testIssuer, Secret: testSecret}).
			Parse(f.Templater.tokens[0], &claims)
		require.NoError(t, err)

		assert.Equal(t, "user@example.com", claims.Data.Email)
		assert.Equal(t, testIssuer, claims.Issuer)
		assert.NotEmpty(t, claims.ID, "token must carry a jti, it is what makes it single-use")
		assert.WithinDuration(t, time.Now().Add(confTokenTtlDuration), claims.ExpiresAt.Time, time.Minute)
	})

	t.Run("should give every login its own token", func(t *testing.T) {
		f := setupAuth(t)

		f.call(t, f.Service.Login, LoginSchema{Email: "user@example.com"})
		f.call(t, f.Service.Login, LoginSchema{Email: "user@example.com"})

		require.Len(t, f.Templater.tokens, 2)
		assert.NotEqual(t, f.Templater.tokens[0], f.Templater.tokens[1])
	})

	t.Run("should return 422 for an invalid email", func(t *testing.T) {
		f := setupAuth(t)

		w := f.call(t, f.Service.Login, LoginSchema{Email: "not-an-email"})
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
		assert.Empty(t, f.Email.sent)
	})

	t.Run("should return 422 for an empty email", func(t *testing.T) {
		f := setupAuth(t)

		w := f.call(t, f.Service.Login, LoginSchema{})
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
		assert.Empty(t, f.Email.sent)
	})

	t.Run("should return 500 when the email cannot be sent", func(t *testing.T) {
		f := setupAuth(t)
		f.Email.err = errors.New("smtp is down")

		w := f.call(t, f.Service.Login, LoginSchema{Email: "user@example.com"})
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return 500 when the email cannot be templated", func(t *testing.T) {
		f := setupAuth(t)
		f.Templater.err = errors.New("broken template")

		w := f.call(t, f.Service.Login, LoginSchema{Email: "user@example.com"})
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Empty(t, f.Email.sent)
	})
}

func Test_AuthService_Confirm(t *testing.T) {
	t.Run("should create the user and open a session for a token from the email", func(t *testing.T) {
		f := setupAuth(t)
		authToken := f.login(t, "user@example.com")

		w := f.call(t, f.Service.Confirm, ConfirmationSchema{Token: authToken})
		assert.Equal(t, http.StatusOK, w.Code)

		user, err := f.Storage.GetUserByEmail(t.Context(), domain.Email("user@example.com"))
		require.NoError(t, err)
		assert.NotEmpty(t, user.Name, "a new user gets a generated name")

		assert.Equal(t, []string{string(user.ID)}, f.Session.setUserIDs)

		var resp struct {
			Data domain.User
		}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, user, resp.Data)
	})

	t.Run("should reject the second use of the same token", func(t *testing.T) {
		f := setupAuth(t)
		authToken := f.login(t, "user@example.com")

		first := f.call(t, f.Service.Confirm, ConfirmationSchema{Token: authToken})
		require.Equal(t, http.StatusOK, first.Code)

		second := f.call(t, f.Service.Confirm, ConfirmationSchema{Token: authToken})
		assert.Equal(t, http.StatusBadRequest, second.Code)
		assert.Len(t, f.Session.setUserIDs, 1, "the replayed token must not open a second session")
	})

	t.Run("should log an existing user in without creating a second account", func(t *testing.T) {
		f := setupAuth(t)

		existing, err := f.Storage.CreateUser(t.Context(), domain.User{
			ID:    domain.NewUserID(),
			Name:  "Existing user",
			Email: domain.Email("user@example.com"),
		})
		require.NoError(t, err)

		w := f.call(t, f.Service.Confirm, ConfirmationSchema{Token: f.login(t, "user@example.com")})
		assert.Equal(t, http.StatusOK, w.Code)

		assert.Equal(t, []string{string(existing.ID)}, f.Session.setUserIDs)

		stored, err := f.Storage.GetUserByEmail(t.Context(), domain.Email("user@example.com"))
		require.NoError(t, err)
		assert.Equal(t, existing, stored, "the existing name must survive the login")
	})

	t.Run("should reject a token signed with another secret", func(t *testing.T) {
		f := setupAuth(t)
		forged := makeToken(t, "another_secret", "user@example.com", time.Now().Add(time.Hour), "jti-1")

		w := f.call(t, f.Service.Confirm, ConfirmationSchema{Token: forged})
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Empty(t, f.Session.setUserIDs)

		_, err := f.Storage.GetUserByEmail(t.Context(), domain.Email("user@example.com"))
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("should reject an expired token", func(t *testing.T) {
		f := setupAuth(t)
		// beyond the one minute of leeway the token parser allows
		expired := makeToken(t, testSecret, "user@example.com", time.Now().Add(-time.Hour), "jti-1")

		w := f.call(t, f.Service.Confirm, ConfirmationSchema{Token: expired})
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Empty(t, f.Session.setUserIDs)
	})

	t.Run("should reject a malformed token", func(t *testing.T) {
		f := setupAuth(t)

		w := f.call(t, f.Service.Confirm, ConfirmationSchema{Token: "not-a-jwt"})
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Empty(t, f.Session.setUserIDs)
	})

	t.Run("should return 500 when the session cannot be set", func(t *testing.T) {
		f := setupAuth(t)
		authToken := f.login(t, "user@example.com")
		f.Session.err = errors.New("no cookie for you")

		w := f.call(t, f.Service.Confirm, ConfirmationSchema{Token: authToken})
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func Test_AuthService_Logout(t *testing.T) {
	t.Run("should clear the session", func(t *testing.T) {
		f := setupAuth(t)

		w := f.call(t, f.Service.Logout, nil)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, 1, f.Session.cleared)
	})
}
