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
		Sections: workout.Sections,
	}
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
	ID   string
	Name string
}

type TrackSchema struct {
	Track   TrackInfo
	IsOwner bool
}
