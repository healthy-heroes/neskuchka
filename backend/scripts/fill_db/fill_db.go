package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"github.com/healthy-heroes/neskuchka/backend/app/store"
	"github.com/healthy-heroes/neskuchka/backend/app/store/datastore"
	"github.com/healthy-heroes/neskuchka/backend/app/store/db"
)

func main() {
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
	database, err := db.NewSqlite(dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Create datastore
	ds := datastore.NewDataStore(database)

	// Initialize tables
	if err := ds.Exercise.InitTables(); err != nil {
		log.Fatalf("Failed to initialize exercise tables: %v", err)
	}
	if err := ds.User.InitTables(); err != nil {
		log.Fatalf("Failed to initialize user tables: %v", err)
	}
	if err := ds.Track.InitTables(); err != nil {
		log.Fatalf("Failed to initialize track tables: %v", err)
	}
	if err := ds.Workout.InitTables(); err != nil {
		log.Fatalf("Failed to initialize workout tables: %v", err)
	}

	// Begin transaction
	tx, err := database.Beginx()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Fatalf("Failed to commit transaction: %v", err)
		}
	}()

	// Create user
	userLogin := "first_user"
	userID := store.UserID("usr_" + userLogin)
	firstUser := &store.User{
		ID:    userID,
		Name:  "First User",
		Email: "first_user@example.com",
	}
	fmt.Printf("Creating user: %+v\n", firstUser)
	if _, err = ds.User.Create(firstUser); err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	// Create track
	trackID := store.TrackID(uuid.New().String())
	sportsTrack := &store.Track{
		ID:      trackID,
		Name:    "Нескучный спорт",
		OwnerID: userID,
	}
	fmt.Printf("Creating track: %+v\n", sportsTrack)
	if _, err = ds.Track.Create(sportsTrack); err != nil {
		log.Fatalf("Failed to create track: %v", err)
	}

	// Create exercises
	exercises := map[string]*store.Exercise{
		"Раскрытия в планке": {
			Name:        "Раскрытия в планке",
			Slug:        store.ExerciseSlug("plank-hip-opening"),
			Description: "",
		},
		"Ягодичные марши": {
			Name:        "Ягодичные марши",
			Slug:        store.ExerciseSlug("glute-march"),
			Description: "",
		},
		"Джампинг джек": {
			Name:        "Джампинг джек",
			Slug:        store.ExerciseSlug("jumping-jack"),
			Description: "",
		},
		"C пола на грудь + 2 выпада назад": {
			Name:        "C пола на грудь + 2 выпада назад",
			Slug:        store.ExerciseSlug("push-up-with-back-drop"),
			Description: "",
		},
		"Становая на одной": {
			Name:        "Становая на одной",
			Slug:        store.ExerciseSlug("deadlift-on-one-leg"),
			Description: "",
		},
		"Пресс на прямых руки над головой": {
			Name:        "Пресс на прямых руки над головой",
			Slug:        store.ExerciseSlug("situps-with-hands-over-head"),
			Description: "",
		},
		"Стол": {
			Name:        "Стол",
			Slug:        store.ExerciseSlug("table"),
			Description: "",
		},
		"Наклоны вперед": {
			Name:        "Наклоны вперед",
			Slug:        store.ExerciseSlug("forward-bend"),
			Description: "",
		},
		"Качающиеся планки": {
			Name:        "Качающиеся планки",
			Slug:        store.ExerciseSlug("plank-with-jumping-jack"),
			Description: "",
		},
		"Становая тяга": {
			Name:        "Становая тяга",
			Slug:        store.ExerciseSlug("deadlift"),
			Description: "",
		},
		"Приседания": {
			Name:        "Приседания",
			Slug:        store.ExerciseSlug("squats"),
			Description: "",
		},
		"С груди над головой": {
			Name:        "С груди над головой",
			Slug:        store.ExerciseSlug("push-up-with-hands-over-head"),
			Description: "",
		},
	}

	for name, ex := range exercises {
		fmt.Printf("Creating exercise: %s (%s)\n", name, ex.Slug)
		if _, err = ds.Exercise.Create(ex); err != nil {
			log.Fatalf("Failed to create exercise: %v", err)
		}
	}

	// Create first workout (Jan 31, 2025)
	firstWorkoutID := store.CreateWorkoutId()
	firstWorkout := &store.Workout{
		ID:      firstWorkoutID,
		Date:    time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC),
		TrackID: trackID,
		Sections: []store.WorkoutSection{
			{
				Title: "Разминка",
				Protocol: store.WorkoutProtocol{
					Type:        store.WorkoutProtocolTypeDefault,
					Title:       "3 раунда",
					Description: "",
				},
				Exercises: []store.WorkoutExercise{
					{
						ExerciseSlug: exercises["Раскрытия в планке"].Slug,
						Repetitions:  20,
						Description:  "20 раскрытий в планке",
					},
					{
						ExerciseSlug: exercises["Ягодичные марши"].Slug,
						Repetitions:  20,
						Description:  "20 ягодичных маршей",
					},
					{
						ExerciseSlug: exercises["Джампинг джек"].Slug,
						Repetitions:  30,
						Description:  "30 джампинг джеков",
					},
				},
			},
			{
				Title: "Комплекс",
				Protocol: store.WorkoutProtocol{
					Type:        store.WorkoutProtocolTypeDefault,
					Title:       "5 раундов",
					Description: "*можно использовать спортивные снаряды",
				},
				Exercises: []store.WorkoutExercise{
					{
						ExerciseSlug: exercises["C пола на грудь + 2 выпада назад"].Slug,
						Repetitions:  10,
						Description:  "10 отжиманий",
					},
					{
						ExerciseSlug:    exercises["Становая на одной"].Slug,
						Repetitions:     20,
						RepetitionsText: "10+10",
						Description:     "20 становых на одной",
					},
					{
						ExerciseSlug: exercises["Пресс на прямых руки над головой"].Slug,
						Repetitions:  10,
						Description:  "10 пресса на прямых руках над головой",
					},
				},
			},
		},
	}

	fmt.Printf("Creating workout: %+v\n", firstWorkout)
	if _, err = ds.Workout.Create(firstWorkout); err != nil {
		log.Fatalf("Failed to create workout: %v", err)
	}

	// Create second workout (Feb 3, 2025)
	secondWorkoutID := store.CreateWorkoutId()
	secondWorkout := &store.Workout{
		ID:      secondWorkoutID,
		Date:    time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC),
		TrackID: trackID,
		Sections: []store.WorkoutSection{
			{
				Title: "Разминка",
				Protocol: store.WorkoutProtocol{
					Type:        store.WorkoutProtocolTypeDefault,
					Title:       "3 раунда",
					Description: "",
				},
				Exercises: []store.WorkoutExercise{
					{
						ExerciseSlug: exercises["Стол"].Slug,
						Repetitions:  10,
						Description:  "10 столов",
					},
					{
						ExerciseSlug: exercises["Наклоны вперед"].Slug,
						Repetitions:  10,
						Description:  "10 наклонов вперед",
					},
					{
						ExerciseSlug: exercises["Качающиеся планки"].Slug,
						Repetitions:  20,
						Description:  "20 качающихся планок",
					},
				},
			},
			{
				Title: "Комплекс",
				Protocol: store.WorkoutProtocol{
					Type:        store.WorkoutProtocolTypeDefault,
					Title:       "5 раундов",
					Description: "",
				},
				Exercises: []store.WorkoutExercise{
					{
						ExerciseSlug: exercises["Становая на одной"].Slug,
						Repetitions:  24,
						Description:  "24 становых на одной",
					},
					{
						ExerciseSlug: exercises["Приседания"].Slug,
						Repetitions:  18,
						Description:  "18 приседаний",
					},
					{
						ExerciseSlug: exercises["С груди над головой"].Slug,
						Repetitions:  12,
						Description:  "12 отжиманий",
					},
				},
			},
		},
	}

	fmt.Printf("Creating workout: %+v\n", secondWorkout)
	if _, err = ds.Workout.Create(secondWorkout); err != nil {
		log.Fatalf("Failed to create workout: %v", err)
	}

	fmt.Println("Database successfully populated!")
}
