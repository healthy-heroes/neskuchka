package cmd

import (
	"context"
	"embed"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/healthy-heroes/neskuchka/backend/app/api"
	"github.com/healthy-heroes/neskuchka/backend/app/store/datastore"
	"github.com/healthy-heroes/neskuchka/backend/app/store/db"
)

//go:embed web
var webFS embed.FS

// ServerCommand is the command for the run server
type ServerCommand struct {
	Address string `long:"address" env:"ADDRESS" default:"127.0.0.1" description:"address"`
	Port    int    `long:"port" env:"PORT" default:"8080" description:"port"`
}

// serverApp holds all active objects
type serverApp struct {
	*ServerCommand

	apiServer *api.Api
	store     *datastore.DataStore
}

func (cmd *ServerCommand) Execute(args []string) error {
	log.Info().Msgf("Starting server on port %d", cmd.Port)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Warn().Msg("Got interrupt signal")
		cancel()
	}()

	app, err := cmd.newServerApp()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create app")
		return err
	}

	if err = app.run(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to terminate app")
		return err
	}

	log.Info().Msg("Server terminated")
	return nil
}

func (cmd *ServerCommand) newServerApp() (*serverApp, error) {
	db, err := db.NewDB(filepath.Join(".", "bin", "app.db"))
	if err != nil {
		return nil, err
	}
	store := datastore.NewDataStore(db)

	apiServer := &api.Api{
		Store: store,
		WebFS: webFS,
	}

	app := &serverApp{
		ServerCommand: cmd,
		apiServer:     apiServer,
		store:         store,
	}

	return app, nil
}

// Run all application objects
func (app *serverApp) run(ctx context.Context) error {
	go func() {
		// shutdown on context cancellation
		<-ctx.Done()

		log.Info().Msg("Handle shutdown...")

		app.apiServer.Shutdown()
	}()

	app.apiServer.Run(app.Address, app.Port)

	return nil
}
