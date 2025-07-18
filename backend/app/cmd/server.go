package cmd

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-pkgz/auth/v2"
	"github.com/go-pkgz/auth/v2/avatar"
	"github.com/go-pkgz/auth/v2/token"
	"github.com/rs/zerolog/log"

	"github.com/healthy-heroes/neskuchka/backend/app/api"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/authproviders"
	"github.com/healthy-heroes/neskuchka/backend/app/store"
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

	authService := cmd.getAuthService(dataStore)

	apiServer := &api.Api{
		Version: cmd.Revision,

		Store:       dataStore,
		AuthService: authService,
		WebFS:       webFS,
	}

	app := &serverApp{
		ServerCommand: cmd,
		apiServer:     apiServer,
		store:         dataStore,
	}

	return app, nil
}

// makeDataStore creates a new data store
func (cmd *ServerCommand) makeDataStore() (*datastore.DataStore, error) {
	log.Info().Msgf("Creating store: %s", cmd.Store.Type)

	switch cmd.Store.Type {
	case db.Sqlite:
		if cmd.Store.Sqlite.Source == "" {
			return nil, fmt.Errorf("sqlite source is not set")
		}

		db, err := db.NewSqlite(cmd.Store.Sqlite.Source)
		if err != nil {
			return nil, fmt.Errorf("failed to create sqlite database: %w", err)
		}

		return datastore.NewDataStore(db), nil

	default:
		return nil, fmt.Errorf("unsupported database type: %s", cmd.Store.Type)
	}
}

// getAuthService creates a new authentication service
func (cmd *ServerCommand) getAuthService(ds *datastore.DataStore) *auth.Service {
	options := auth.Opts{
		//todo: getting secret from server command opts
		SecretReader: token.SecretFunc(func(id string) (string, error) {
			// secret key for JWT
			return "secret", nil
		}),
		ClaimsUpd: token.ClaimsUpdFunc(func(claims token.Claims) token.Claims {
			if claims.User == nil {
				return claims
			}

			user, err := ds.User.Get(store.UserID(claims.User.ID))
			if err != nil {
				log.Error().Err(err).Msgf("Error finding user %s", claims.User.ID)

				claims.User = nil
				return claims
			}

			// update claims with actual user data
			claims.User.Name = user.Name

			return claims
		}),
		TokenDuration:  time.Minute * 5, // token expires in 5 minutes
		CookieDuration: time.Hour * 24,  // cookie expires in 1 day and will enforce re-login
		Issuer:         "neskuchka",
		AvatarStore:    avatar.NewLocalFS("/tmp"),
		Validator: token.ValidatorFunc(func(_ string, claims token.Claims) bool {
			if claims.User == nil {
				log.Error().Msgf("User nil in validator %s", claims)

				return false
			}

			return true
		}),
	}

	service := auth.NewService(options)
	providers := authproviders.NewService(options)

	// todo: make normal email sender
	emailSender := AuthEmailSender{}
	msgTemplate := "Confirmation email, token:  http://localhost:8081/auth/email/login?token={{.Token}}"

	verify := providers.NewVerifyProvider("email", msgTemplate, emailSender,
		func(name string, email string, r *http.Request) (string, error) {
			user, err := ds.User.FindOrCreate(email, name)
			if err != nil && err != store.ErrNotFound {
				log.Error().Err(err).Msgf("Error getting user %s", email)

				return "", err
			}

			return string(user.ID), nil
		})
	service.AddCustomHandler(verify)

	return service
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
