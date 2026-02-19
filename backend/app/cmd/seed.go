package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
	"github.com/healthy-heroes/neskuchka/backend/app/storage/database"
)

type SeedCommand struct {
	Store StoreOptions `group:"store" namespace:"store" env-namespace:"STORE"`

	CommonOptions
}

type SeedRunner struct {
	dataStorage *database.DataStorage
}

func (cmd *SeedCommand) Execute(args []string) error {
	log.Info().Msg("running seed...")

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Warn().Msg("got interrupt signal")
		cancel()
	}()

	runner, err := cmd.createRunner()
	if err != nil {
		return fmt.Errorf("failed create seed runner: %w", err)
	}

	if err := runner.Run(ctx); err != nil {
		return fmt.Errorf("failed execute runner: %w", err)
	}

	log.Info().Msg("database seeded successfully")
	return nil
}

func (cmd *SeedCommand) createRunner() (*SeedRunner, error) {
	log.Info().Msg("creating store")
	log.Info().Msgf("database url: %s", cmd.Store.DB)

	engine, err := database.NewEngine(cmd.Store.DB, database.Opts{Logger: log.Logger})
	if err != nil {
		return nil, fmt.Errorf("failed to create engine: %w", err)
	}

	return &SeedRunner{
		dataStorage: database.NewDataStorage(engine, log.Logger),
	}, nil
}

func (r *SeedRunner) Run(ctx context.Context) error {
	log.Info().Msg("start seed runner")

	go func() {
		// shutdown on context cancellation
		<-ctx.Done()
		log.Info().Msg("runner shutdown...")
		r.dataStorage.Close()
	}()

	admin := domain.User{
		ID:    domain.NewUserID(),
		Name:  "Admin",
		Email: "admin@example.com",
	}

	// check existing admin
	_, err := r.dataStorage.GetUserByEmail(ctx, admin.Email)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return fmt.Errorf("failed to check existing data: %w", err)
	}
	if err == nil {
		log.Info().Msg("data already exists, skipping")
		return nil
	}

	// Create user
	_, err = r.dataStorage.CreateUser(ctx, admin)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Create track
	mainTrack := domain.Track{
		ID:          domain.NewTrackID(),
		Name:        "Нескучный спорт",
		Slug:        domain.TrackSlug("main"),
		Description: "Тренируйтесь с нами — где бы вы ни находились!\nИдеальная программа, чтобы поддерживать форму дома, без специального оборудования.",
		OwnerID:     admin.ID,
	}
	_, err = r.dataStorage.CreateTrack(ctx, mainTrack)
	if err != nil {
		return fmt.Errorf("failed to create track: %w", err)
	}

	// Create first workout (Jan 31, 2025)
	_, err = r.dataStorage.CreateWorkout(ctx, domain.Workout{
		ID:      domain.NewWorkoutID(),
		Date:    time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC),
		TrackID: mainTrack.ID,
		Notes:   "First workout",
		Sections: []domain.WorkoutSection{
			{
				Title: "Разминка",
				Protocol: domain.Protocol{
					Type:        domain.ProtocolTypeCustom,
					Title:       "3 раунда",
					Description: "",
				},
				Exercises: []domain.WorkoutExercise{
					{
						ExerciseSlug: domain.ExerciseSlug("plank-hip-opening"),
						Description:  "20 раскрытий в планке",
					},
				},
			},
			{
				Title: "Комплекс",
				Protocol: domain.Protocol{
					Type:        domain.ProtocolTypeCustom,
					Title:       "5 раундов",
					Description: "*можно использовать спортивные снаряды",
				},
				Exercises: []domain.WorkoutExercise{
					{
						ExerciseSlug: domain.ExerciseSlug("push-up-with-back-drop"),
						Description:  "10 отжиманий",
					},
					{
						ExerciseSlug: domain.ExerciseSlug("deadlift-on-one-leg"),
						Description:  "20 становых на одной",
					},
					{
						ExerciseSlug: domain.ExerciseSlug("situps-with-hands-over-head"),
						Description:  "10 пресса на прямых руках над головой",
					},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create workout: %w", err)
	}

	// Create second workout (Feb 3, 2025)
	_, err = r.dataStorage.CreateWorkout(ctx, domain.Workout{
		ID:      domain.NewWorkoutID(),
		Date:    time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC),
		TrackID: mainTrack.ID,
		Notes:   "Second workout",
		Sections: []domain.WorkoutSection{
			{
				Title: "Разминка",
				Protocol: domain.Protocol{
					Type:        domain.ProtocolTypeCustom,
					Title:       "3 раунда",
					Description: "",
				},
				Exercises: []domain.WorkoutExercise{
					{
						ExerciseSlug: domain.ExerciseSlug("table"),
						Description:  "10 столов",
					},
					{
						ExerciseSlug: domain.ExerciseSlug("forward-bend"),
						Description:  "10 наклонов вперед",
					},
					{
						ExerciseSlug: domain.ExerciseSlug("plank-with-jumping-jack"),
						Description:  "20 качающихся планок",
					},
				},
			},
			{
				Title: "Комплекс",
				Protocol: domain.Protocol{
					Type:        domain.ProtocolTypeCustom,
					Title:       "5 раундов",
					Description: "",
				},
				Exercises: []domain.WorkoutExercise{
					{
						ExerciseSlug: domain.ExerciseSlug("deadlift-on-one-leg"),
						Description:  "24 становых на одной",
					},
					{
						ExerciseSlug: domain.ExerciseSlug("squats"),
						Description:  "18 приседаний",
					},
					{
						ExerciseSlug: domain.ExerciseSlug("push-up-with-hands-over-head"),
						Description:  "12 отжиманий",
					},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create workout: %w", err)
	}

	return nil
}
