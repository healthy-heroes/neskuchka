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
	"github.com/healthy-heroes/neskuchka/backend/app/internal/email"
	"github.com/healthy-heroes/neskuchka/backend/app/storage/database"
)

//go:embed web
var webFS embed.FS

// ServerCommand is the command for the run server
type ServerCommand struct {
	BaseURL string `long:"base-url" env:"BASE_URL" default:"http://localhost:8080" description:"base URL for the application"`

	Address string `long:"address" env:"ADDRESS" default:"127.0.0.1" description:"address"`
	Port    int    `long:"port" env:"PORT" default:"8080" description:"port"`

	Store StoreOptions `group:"store" namespace:"store" env-namespace:"STORE"`
	SMTP  SMTPOptions  `group:"smtp" namespace:"smtp" env-namespace:"SMTP"`

	Secret string `long:"secret" env:"SECRET" description:"secret key for JWT"`

	CommonOptions
}

// StoreOptions defines options for the storage
type StoreOptions struct {
	DB string `long:"db" env:"DB" description:"database URL (sqlite file)"`
}

type SMTPOptions struct {
	Host string `long:"host" env:"HOST" description:"SMTP host"`
	Port int    `long:"port" env:"PORT" description:"SMTP port"`
}

// serverApp holds all active objects
type serverApp struct {
	*ServerCommand

	apiServer *api.Api
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
	dataStore, err := cmd.makeDataStore()
	if err != nil {
		return nil, fmt.Errorf("failed to create data store: %w", err)
	}

	apiServer := &api.Api{
		Version: cmd.Revision,
		Secret:  cmd.Secret,

		DataStore: dataStore,
		WebFS:     webFS,

		EmailTemplater: email.NewTemplate(cmd.BaseURL),
		EmailService: email.NewService(email.Opts{
			Host:   cmd.SMTP.Host,
			Port:   cmd.SMTP.Port,
			Logger: log.Logger,
		}),
	}

	app := &serverApp{
		ServerCommand: cmd,
		apiServer:     apiServer,
		dataStore:     dataStore,
	}

	return app, nil
}

// makeDataStore creates a new data store
func (cmd *ServerCommand) makeDataStore() (*domain.Store, error) {
	log.Info().Msg("creating store")
	log.Info().Msgf("database url: %s", cmd.Store.DB)

	engine, err := database.NewEngine(cmd.Store.DB, database.Opts{Logger: log.Logger})
	if err != nil {
		return nil, fmt.Errorf("failed to create engine: %w", err)
	}

	return domain.NewStore(domain.Opts{
		DataStorage: database.NewDataStorage(engine, log.Logger),
	}), nil
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
