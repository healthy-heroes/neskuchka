package cmd

import (
	"context"
	"embed"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/healthy-heroes/neskuchka/backend/app/api"
	"github.com/healthy-heroes/neskuchka/backend/app/domain"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/storage/database"
	"github.com/healthy-heroes/neskuchka/backend/app/store/datastore"
	"github.com/healthy-heroes/neskuchka/backend/app/store/db"
)

//go:embed web
var webFS embed.FS

// ServerCommand is the command for the run server
type ServerCommand struct {
	Store StoreOptions `group:"store" namespace:"store" env-namespace:"STORE"`

	Address string `long:"address" env:"ADDRESS" default:"127.0.0.1" description:"address"`
	Port    int    `long:"port" env:"PORT" default:"8080" description:"port"`

	Secret string `long:"secret" env:"SECRET" description:"secret key for JWT"`

	CommonOptions
}

// StoreOptions defines options for the storage
type StoreOptions struct {
	Type   db.Type `long:"type" env:"TYPE" description:"type of storage" default:"sqlite"`
	Sqlite struct {
		Source string `long:"source" env:"SOURCE" description:"file name or :memory:"`
	} `group:"sqlite" namespace:"sqlite" env-namespace:"SQLITE"`
}

// serverApp holds all active objects
type serverApp struct {
	*ServerCommand

	apiServer *api.Api
	store     *datastore.DataStore
	dataStore *domain.Store

	CommonOptions
}

// Execute starts the server
func (cmd *ServerCommand) Execute(args []string) error {
	log.Info().Msgf("Starting server on port %d, (revision: %s)", cmd.Port, cmd.Revision)

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
	dataStore, newDataStore, err := cmd.makeDataStore()
	if err != nil {
		return nil, fmt.Errorf("failed to create data store: %w", err)
	}

	apiServer := &api.Api{
		Version: cmd.Revision,
		Secret:  cmd.Secret,

		DataStore: newDataStore,
		Store:     dataStore,
		WebFS:     webFS,
	}

	app := &serverApp{
		ServerCommand: cmd,
		apiServer:     apiServer,
		store:         dataStore,
		dataStore:     newDataStore,
	}

	return app, nil
}

// makeDataStore creates a new data store
func (cmd *ServerCommand) makeDataStore() (*datastore.DataStore, *domain.Store, error) {
	log.Info().Msgf("Creating store: %s", cmd.Store.Type)

	switch cmd.Store.Type {
	case db.Sqlite:
		if cmd.Store.Sqlite.Source == "" {
			return nil, nil, fmt.Errorf("sqlite source is not set")
		}

		// todo remove
		db, err := db.NewSqlite(cmd.Store.Sqlite.Source)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create sqlite database: %w", err)
		}

		engine, err := database.NewSqliteEngine(cmd.Store.Sqlite.Source, log.Logger)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create sqlite engine: %w", err)
		}

		oldDataStore := datastore.NewDataStore(db)
		newDataStore := domain.NewStore(domain.Opts{
			DataStorage: database.NewDataStorage(engine, log.Logger),
		})

		return oldDataStore, newDataStore, nil

	default:
		return nil, nil, fmt.Errorf("unsupported database type: %s", cmd.Store.Type)
	}
}

// fake email sender
type AuthEmailSender struct{}

func (s AuthEmailSender) Send(email string, text string) error {
	log.Info().Msgf("Sending email to %s:\n\n%s\n", email, text)
	return nil
}

// run starts all application objects
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
