package api

import (
	"context"
	"fmt"
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
	httpServer *http.Server
	lock       sync.Mutex

	store  *datastore.DataStore
	public *public_api.PublicAPI
}

func NewApi(store *datastore.DataStore) *Api {
	return &Api{
		store: store,
	}
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

	api.public = public_api.NewPublicAPI(api.store)

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

	return router
}
