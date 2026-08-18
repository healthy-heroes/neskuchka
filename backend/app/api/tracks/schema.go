package tracks

import (
	"time"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
)

// WorkoutInfo is a workout info for reading and creating workouts
type WorkoutInfo struct {
	ID       string
	TrackID  string
	Date     string
	Notes    string
	Sections []domain.WorkoutSection
}

func MakeWorkoutInfo(workout domain.Workout) WorkoutInfo {
	return WorkoutInfo{
		ID:       string(workout.ID),
		TrackID:  string(workout.TrackID),
		Date:     workout.Date.Format(time.DateOnly),
		Notes:    workout.Notes,
		Sections: makeSections(workout.Sections),
	}
}

// makeSections holds the API contract that every list is an array, never null.
//
// Both cases happen for real: a skill section is often just a heading with no
// exercises, and an exercise in a warm-up often has no prescription. A nil slice
// marshals to null, and clients iterate these lists without checking.
func makeSections(sections []domain.WorkoutSection) []domain.WorkoutSection {
	result := make([]domain.WorkoutSection, 0, len(sections))

	for _, section := range sections {
		exercises := make([]domain.WorkoutExercise, 0, len(section.Exercises))
		for _, exercise := range section.Exercises {
			if exercise.Prescription == nil {
				exercise.Prescription = []string{}
			}
			exercises = append(exercises, exercise)
		}

		section.Exercises = exercises
		result = append(result, section)
	}

	return result
}

func (w *WorkoutInfo) toDomain() (domain.Workout, error) {
	date, err := time.Parse(time.DateOnly, w.Date)
	if err != nil {
		return domain.Workout{}, err
	}

	return domain.Workout{
		ID:       domain.WorkoutID(w.ID),
		TrackID:  domain.TrackID(w.TrackID),
		Date:     date,
		Notes:    w.Notes,
		Sections: w.Sections,
	}, nil
}

func MakeWorkoutInfos(workouts []domain.Workout) []WorkoutInfo {
	workoutInfos := make([]WorkoutInfo, 0, len(workouts))
	for _, workout := range workouts {
		workoutInfos = append(workoutInfos, MakeWorkoutInfo(workout))
	}
	return workoutInfos
}

type WorkoutSchema struct {
	Workout WorkoutInfo
}

type WorkoutsSchema struct {
	Workouts []WorkoutInfo
}

type TrackInfo struct {
	ID          string
	Name        string
	Description string

	Author AuthorInfo
}

// AuthorInfo is the track owner shown next to the track description.
// No avatar yet: the user service only emits an avatar URL when the file
// actually exists, and wiring the avatar store in here is not worth it.
type AuthorInfo struct {
	ID   string
	Name string
}

func MakeTrackInfo(track domain.Track, author domain.User) TrackInfo {
	return TrackInfo{
		ID:          string(track.ID),
		Name:        track.Name,
		Description: track.Description,
		Author: AuthorInfo{
			ID:   string(author.ID),
			Name: author.Name,
		},
	}
}

type TrackSchema struct {
	Track   TrackInfo
	IsOwner bool
}
