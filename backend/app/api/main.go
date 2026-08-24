package api

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	chiMW "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/rs/zerolog/log"

	"github.com/healthy-heroes/neskuchka/backend/app/api/auth"
	"github.com/healthy-heroes/neskuchka/backend/app/api/httpx"
	mw "github.com/healthy-heroes/neskuchka/backend/app/api/middlewares"
	"github.com/healthy-heroes/neskuchka/backend/app/api/tracks"
	api_user "github.com/healthy-heroes/neskuchka/backend/app/api/user"
	"github.com/healthy-heroes/neskuchka/backend/app/domain"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/email"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/session"
	"github.com/healthy-heroes/neskuchka/backend/app/storage/avatarstorage"
)

const Issuer = "Neskuchka"
const prefixApi = "/api/v1"

// requestTimeout is how long a request may take before chi cancels its
// context. Anything a route makes a client wait for has to fit inside it.
const requestTimeout = 10 * time.Second

// headerTimeout is how long a client has to send its headers. They are small
// and every client sends them at once, so this needs no room to breathe.
const headerTimeout = 5 * time.Second

// maxRequestTimeout is the ceiling every API route sits under, not the budget
// any of them is meant to use: routes shorten it to what they actually need,
// and one that forgets to is left bounded rather than running forever. It stays
// above the longest route on purpose, so that the route's own timeout is the
// one that fires and reaches whoever is waiting.
const maxRequestTimeout = 90 * time.Second

// clientIPKey keys the rate limiters off the client IP resolved by chi's
// ClientIPFrom* middleware. CanonicalizeIP buckets IPv6 clients by their /64,
// so one client can't get a fresh budget per address in its prefix.
//
// Today the resolver is ClientIPFromRemoteAddr — the TCP peer. That is correct
// only while the server is reached directly. Behind a reverse proxy every
// client collapses into the proxy's single bucket; that case needs
// ClientIPFromXFF with the proxy's CIDR, which in turn needs the deployment
// topology to be known.
func clientIPKey(r *http.Request) (string, error) {
	return httprate.CanonicalizeIP(chiMW.GetClientIP(r.Context())), nil
}

// Api is an API server
type Api struct {
	Version string
	Secret  string

	DataStore   *domain.Store
	AvatarStore *avatarstorage.Storage
	WebFS       fs.FS

	httpServer *http.Server
	lock       sync.Mutex

	EmailTemplater *email.Templater
	EmailService   *email.Service
}

// Run the listener and request's router, starts the API server
func (api *Api) Run(address string, port int) {
	api.lock.Lock()
	api.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", address, port),
		Handler: api.Handler(),

		// Without it a client holds a connection open on any route, and without
		// logging in, simply by never finishing its headers.
		ReadHeaderTimeout: headerTimeout,
	}
	api.lock.Unlock()

	err := api.httpServer.ListenAndServe()
	log.Warn().Err(err).Msg("Api server terminated")
}

// Shutdown shuts down the API server
func (api *Api) Shutdown() {
	log.Info().Msg("Shutting down api server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	api.lock.Lock()

	if api.httpServer != nil {
		if err := api.httpServer.Shutdown(ctx); err != nil {
			log.Warn().Err(err).Msg("Api http server shutdown error")
		}
		log.Info().Msg("Api http server shutdown completed")
	}

	api.lock.Unlock()
}

// routes is setting up routes for the API
func (api *Api) Handler() *chi.Mux {
	router := chi.NewRouter()
	session := session.NewManager(session.Opts{
		Logger: log.Logger,
		Issuer: Issuer,
		Secret: api.Secret,
	})

	// common middlewares
	router.Use(chiMW.Logger)
	// resolves the client IP the rate limiters key off, see clientIPKey
	router.Use(chiMW.ClientIPFromRemoteAddr)
	router.Use(session.Verifier())

	// CORS middleware
	corsMw := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-XSRF-Token", "X-JWT"},
		ExposedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	router.Use(corsMw.Handler)

	// ping route
	router.With(
		httprate.LimitBy(600, time.Minute, clientIPKey),
	).Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	// api routes. Each group states its own timeout, because chi's Timeout can
	// only ever shorten — context.WithTimeout keeps the earlier of the two
	// deadlines — and the avatar upload needs longer than the rest.
	router.Route(prefixApi, func(r chi.Router) {
		r.Use(httprate.LimitBy(60, time.Minute, clientIPKey))
		r.Use(chiMW.Timeout(maxRequestTimeout))

		api.addAuthRoutes(r, session)
		api.addUserRoutes(r, session)
		api.addTracksRoutes(r, session)
	})

	api.addStaticRoutes(router)

	return router
}

// addAuthRoutes is adding auth routes
func (api *Api) addAuthRoutes(router chi.Router, session *session.Manager) {
	h := auth.NewService(api.DataStore, session, auth.Opts{
		Issuer: "Neskuchka",
		Secret: api.Secret,
		Logger: log.Logger,

		EmailSender:    api.EmailService,
		EmailTemplater: api.EmailTemplater,
	})

	router.Route("/auth", func(r chi.Router) {
		r.Use(chiMW.Timeout(requestTimeout))

		r.Post("/login", h.Login)
		r.Post("/login/confirm", h.Confirm)
		r.Post("/logout", h.Logout)
	})
}

// Decoding a picture the user chose is the one thing here that costs the
// process hundreds of megabytes, and httprate counts per client — N clients
// buy N decodes at once. Hence a ceiling of its own; two of them at their worst
// is a little over 200mb, and the wait has to leave the decode room in
// requestTimeout.
const (
	maxAvatarUploads  = 2
	avatarUploadQueue = 4
	avatarUploadWait  = 5 * time.Second

	// A slot is held for as long as the handler runs, and reading the body is
	// part of that, so an upload that stops arriving has to be cut off — two
	// clients trickling their bytes would otherwise close avatar uploads for
	// everybody. Eight megabytes, the size cap in api/user, inside a minute
	// asks about a megabit per second of whoever is uploading.
	maxAvatarUploadTime = time.Minute

	// avatarDecodeTime is what is left over for turning the bytes into a stored
	// avatar once they have all arrived — decoding, scaling, one write.
	avatarDecodeTime = 10 * time.Second

	// avatarUploadTimeout is the budget for all of it, and it is a sum rather
	// than a number so that raising any part above raises the whole. Getting
	// this wrong is quiet and expensive: the body would still be read in full
	// and decoded, only for the write to find a context that expired while it
	// was happening.
	avatarUploadTimeout = avatarUploadWait + maxAvatarUploadTime + avatarDecodeTime
)

// addUserRoutes is adding user routes
func (api *Api) addUserRoutes(router chi.Router, session *session.Manager) {
	avatarURLFunc := func(userID domain.UserID) string {
		return fmt.Sprintf("%s/user/%s/avatar", prefixApi, string(userID))
	}

	h := api_user.NewService(api.DataStore, api_user.Opts{
		Logger:        log.Logger,
		AvatarStore:   api.AvatarStore,
		AvatarURLFunc: avatarURLFunc,
	})

	router.Route("/user", func(r chi.Router) {
		r.Route("/me", func(r chi.Router) {
			r.Use(session.Authenticator(httpx.RenderUnauthorized))

			// An upload is the one request here that a slow connection stretches
			// for honest reasons, so it takes avatarUploadTimeout instead of the
			// budget the routes below share. The throttle goes inside it and the
			// read deadline inside the throttle, so that the deadline starts
			// counting when the upload gets its slot rather than while it waits.
			r.With(
				chiMW.Timeout(avatarUploadTimeout),
				mw.Throttle(maxAvatarUploads, avatarUploadQueue, avatarUploadWait),
				mw.ReadDeadline(log.Logger, maxAvatarUploadTime),
			).Post("/avatar", h.UploadAvatar)

			r.Group(func(r chi.Router) {
				r.Use(chiMW.Timeout(requestTimeout))

				r.Get("/", h.Me)

				r.Get("/avatar", h.MyAvatar)
				r.Delete("/avatar", h.DeleteAvatar)

				r.Get("/settings", h.GetSettings)
				r.Put("/settings", h.UpdateSettings)
			})
		})

		r.Route("/{id}", func(r chi.Router) {
			r.Use(chiMW.Timeout(requestTimeout))

			r.Get("/avatar", h.UserAvatar)
		})
	})
}

// addTracksRoutes is adding tracks routes
// temporary working with concrete main track routes
func (api *Api) addTracksRoutes(router chi.Router, session *session.Manager) {
	h := tracks.NewService(api.DataStore, tracks.Opts{
		Logger: log.Logger,
	})

	auth := session.Authenticator(httpx.RenderUnauthorized)

	router.Route("/tracks/main", func(r chi.Router) {
		r.Use(chiMW.Timeout(requestTimeout))

		r.Get("/", h.GetMainTrack)
		r.Get("/last_workouts", h.GetMainTrackLastWorkouts)

		r.Get("/workouts/{id}", h.GetWorkout)

		r.With(auth).Put("/", h.UpdateMainTrack)

		// The whole track, drafts included — owner only, checked in the domain
		r.With(auth).Get("/workouts", h.GetMainTrackWorkouts)

		r.With(auth).Post("/workouts", h.CreateWorkout)
		r.With(auth).Put("/workouts/{id}", h.UpdateWorkout)
		r.With(auth).Delete("/workouts/{id}", h.DeleteWorkout)
	})
}

// addStaticRoutes is adding static routes
func (api *Api) addStaticRoutes(router *chi.Mux) {
	indexHTML, err := fs.ReadFile(api.WebFS, "web/index.html")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read index.html")
	}

	staticFS, _ := fs.Sub(api.WebFS, "web")

	router.Route("/", func(r chi.Router) {
		r.Use(httprate.LimitBy(60, time.Minute, clientIPKey))
		r.Use(chiMW.Timeout(requestTimeout))
		r.Use(mw.CacheControl(10*time.Minute, api.Version))

		r.Handle("/favicon.*", http.FileServer(http.FS(staticFS)))
		r.Handle("/assets/*", http.FileServer(http.FS(staticFS)))
		r.Handle("/img/*", http.FileServer(http.FS(staticFS)))

		//todo: Подумать как улучшить
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write(indexHTML)
		})
	})
}
