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
	"github.com/healthy-heroes/neskuchka/backend/app/domain"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/email"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/session"
)

const Issuer = "Neskuchka"

// Api is an API server
type Api struct {
	Version string
	Secret  string

	DataStore *domain.Store
	WebFS     fs.FS

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
		httprate.LimitByIP(600, time.Minute),
	).Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	// api routes
	router.Route("/api/v1", func(r chi.Router) {
		r.Use(httprate.LimitByIP(60, time.Minute))
		r.Use(chiMW.Timeout(10 * time.Second))

		api.addAuthRoutes(r, session)
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
		r.Post("/login", h.Login)
		r.Post("/login/confirm", h.Confirm)
		r.Post("/logout", h.Logout)

		r.With(session.Authenticator(httpx.RenderUnauthorized)).Get("/user", h.User)
	})
}

// addTracksRoutes is adding tracks routes
// temporary working with concrete main track routes
func (api *Api) addTracksRoutes(router chi.Router, session *session.Manager) {
	h := tracks.NewService(api.DataStore, session, tracks.Opts{
		Logger: log.Logger,
	})

	auth := session.Authenticator(httpx.RenderUnauthorized)

	router.Route("/tracks/main", func(r chi.Router) {
		r.Get("/", h.GetMainTrack)
		r.Get("/last_workouts", h.GetMainTrackLastWorkouts)

		r.Get("/workouts/{id}", h.GetWorkout)

		r.With(auth).Post("/workouts", h.CreateWorkout)
		r.With(auth).Put("/workouts/{id}", h.UpdateWorkout)
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
		r.Use(httprate.LimitByIP(60, time.Minute))
		r.Use(chiMW.Timeout(10 * time.Second))
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
