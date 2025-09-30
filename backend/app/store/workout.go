package store

import (
	"github.com/oklog/ulid/v2"
)

type WorkoutID string
type WorkoutProtocolType string

const (
	WorkoutProtocolTypeDefault WorkoutProtocolType = "default"
)

type Workout struct {
	ID WorkoutID

	Date    string
	TrackID TrackID

	Sections []WorkoutSection
}

type WorkoutSection struct {
	Title     string
	Protocol  WorkoutProtocol
	Exercises []WorkoutExercise
}

type WorkoutProtocol struct {
	Type        WorkoutProtocolType
	Title       string
	Description string
}

type WorkoutExercise struct {
	ExerciseSlug ExerciseSlug

	Repetitions     int
	RepetitionsText string

	Weight     int
	WeightText string

	Description string
}

func CreateWorkoutId() WorkoutID {
	id := ulid.Make()
	return WorkoutID(id.String())
}

type WorkoutStore interface {
	Store

	Create(workout *Workout) (*Workout, error)
	Update(workout *Workout) (*Workout, error)
	Get(id WorkoutID) (*Workout, error)
	Find(criteria *WorkoutFindCriteria) ([]*Workout, error)
}

type WorkoutFindCriteria struct {
	TrackID TrackID
	Limit   int
}

func ExtractSlugsFromWorkouts(workouts []*Workout) []ExerciseSlug {
	slugs := make(map[ExerciseSlug]bool)

	for _, workout := range workouts {
		for _, section := range workout.Sections {
			for _, exercise := range section.Exercises {
				slugs[exercise.ExerciseSlug] = true
			}
		}
	}

	slugsList := make([]ExerciseSlug, 0)
	for slug := range slugs {
		slugsList = append(slugsList, slug)
	}

	return slugsList
}

func ClearUnknownExercisesSlugs(workout *Workout, knownSlugs map[ExerciseSlug]bool) {
	for i := range workout.Sections {
		for j := range workout.Sections[i].Exercises {
			if _, ok := knownSlugs[workout.Sections[i].Exercises[j].ExerciseSlug]; !ok {
				workout.Sections[i].Exercises[j].ExerciseSlug = ""
			}
		}
	}
}
