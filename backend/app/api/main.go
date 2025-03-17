package api

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

type Api struct {
	httpServer *http.Server

	lock sync.Mutex

	public *PublicMethods
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

	api.public = &PublicMethods{}

	router.Use(middleware.Logger)

	// make mw
	router.Get("/api/v1/ping", api.public.pingCtrl)
	return router
}
