package api

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"

	"github.com/healthy-heroes/neskuchka/backend/app/api/public_api"
	"github.com/healthy-heroes/neskuchka/backend/app/store/datastore"
)

type Api struct {
	Store *datastore.DataStore
	WebFS embed.FS

	httpServer *http.Server
	lock       sync.Mutex

	public *public_api.PublicAPI
}

func (api *Api) Run(address string, port int) {
	api.lock.Lock()
	api.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", address, port),
		Handler: api.routes(),
	}
	api.lock.Unlock()

	err := api.httpServer.ListenAndServe()
	log.Warn().Err(err).Msg("Api server terminated")
}

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

func (api *Api) routes() *chi.Mux {
	router := chi.NewRouter()

	api.public = public_api.NewPublicAPI(api.Store)

	// middlewares
	router.Use(middleware.Logger)

	// CORS middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-XSRF-Token", "X-JWT"},
		ExposedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	router.Use(corsMiddleware.Handler)

	router.Route("/api/v1", func(r chi.Router) {
		api.public.InitRoutes(r)
	})

	api.addStaticRoutes(router)

	return router
}

func (api *Api) addStaticRoutes(router *chi.Mux) {
	indexHTML, err := api.WebFS.ReadFile("web/index.html")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read index.html")
	}

	staticFS, _ := fs.Sub(api.WebFS, "web")

	router.Route("/", func(r chi.Router) {
		r.Handle("/favicon.*", http.FileServer(http.FS(staticFS)))
		r.Handle("/assets/*", http.FileServer(http.FS(staticFS)))
		r.Handle("/img/*", http.FileServer(http.FS(staticFS)))

		//todo: Подумать как улучшить
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			if checkWebPath(r.URL.Path) {
				w.WriteHeader(http.StatusOK)
				w.Write(indexHTML)
				return
			}

			http.NotFound(w, r)
		})
	})
}

func checkWebPath(path string) bool {
	switch {
	case path == "/":
		return true
	case path == "/welcome":
		return true
	default:
		return false
	}
}
