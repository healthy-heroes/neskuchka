package domain

import "time"

type WorkoutID string
type WorkoutSlug string
type WorkoutStatus string

const (
	WorkoutStatusPublished = WorkoutStatus("published")
)

type Workout struct {
	ID      WorkoutID
	Slug    WorkoutSlug
	TrackID TrackID

	Date   time.Time
	Status WorkoutStatus
	Notes  string

	Sections []WorkoutSection

	PublishedAt time.Time
}

type WorkoutSection struct {
	Title     string
	Protocol  Protocol
	Exercises []WorkoutExercise
}

type WorkoutExercise struct {
	ExerciseSlug ExerciseSlug

	Description string
}
