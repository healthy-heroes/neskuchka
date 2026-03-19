package devcmd

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/healthy-heroes/neskuchka/backend/app/cmd"
	"github.com/healthy-heroes/neskuchka/backend/app/domain"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/session"
	"github.com/healthy-heroes/neskuchka/backend/app/storage/datastorage"
	"github.com/healthy-heroes/neskuchka/backend/app/storage/db"
)

type TokenCommand struct {
	Store cmd.StoreOptions `group:"store" namespace:"store" env-namespace:"STORE"`

	Email string `required:"true" long:"email" env:"EMAIL" description:"email"`

	cmd.CommonOptions
}

func (cmd *TokenCommand) Execute(args []string) error {
	log.Info().Msg("[dev:token] running...")

	ctx := context.Background()

	engine, err := db.NewEngine(cmd.Store.DB, db.Opts{Logger: log.Logger})
	if err != nil {
		log.Error().Err(err).Msg("failed to create engine")
		return err
	}
	defer engine.Close()

	store := domain.NewStore(domain.Opts{Storage: datastorage.New(engine, log.Logger)})

	user, err := store.FindOrCreateUser(ctx, domain.User{Email: domain.Email(cmd.Email)})
	if err != nil {
		log.Error().Err(err).Msg("failed to find or create user")
		return err
	}

	sm := session.NewManager(session.Opts{
		Issuer:          "DevIssuer",
		Secret:          "__SECRET__",
		SessionDuration: 30 * 24 * time.Hour,
	})

	token, err := sm.Token(string(user.ID))
	if err != nil {
		log.Error().Err(err).Msg("failed to generate token")
		return err
	}

	log.Info().Msgf("token: %s", token)

	return nil
}
