package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
	"github.com/healthy-heroes/neskuchka/backend/app/storage/database"
)

func main() {
	logger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel)

	// Database file path
	dbPath := "./bin/app.db"

	// Make sure directory exists
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	// Delete database file if it exists
	if _, err := os.Stat(dbPath); err == nil {
		if err := os.Remove(dbPath); err != nil {
			log.Fatalf("Failed to delete existing database file: %v", err)
		}
		fmt.Println("Existing database file deleted")
	}

	// Initialize DB
	engine, err := database.NewSqliteEngine(dbPath, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer engine.Close()

	dataStorage := database.NewDataStorage(engine, logger)

	// Begin transaction
	ctx := context.Background()
	tx, err := engine.Beginx()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to begin transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to commit transaction")
		}
	}()

	// Create user
	admin := domain.User{
		ID:    domain.NewUserID(),
		Name:  "Admin",
		Email: "admin@example.com",
	}
	_, err = dataStorage.CreateUser(ctx, admin)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create user")
	}

	// Create track
	mainTrack := domain.Track{
		ID:          domain.NewTrackID(),
		Name:        "Нескучный спорт",
		Slug:        domain.TrackSlug("main"),
		Description: "Тренируйтесь с нами — где бы вы ни находились!\nИдеальная программа, чтобы поддерживать форму дома, без специального оборудования.",
		OwnerID:     admin.ID,
	}
	_, err = dataStorage.CreateTrack(ctx, mainTrack)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create track")
	}

	// Create first workout (Jan 31, 2025)
	_, err = dataStorage.CreateWorkout(ctx, domain.Workout{
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
		logger.Fatal().Err(err).Msg("Failed to create workout")
	}

	// Create second workout (Feb 3, 2025)
	_, err = dataStorage.CreateWorkout(ctx, domain.Workout{
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
		logger.Fatal().Err(err).Msg("Failed to create workout")
	}

	logger.Info().Msg("Database successfully populated!")
}
