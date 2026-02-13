package tracks

import (
	"time"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
)

type WorkoutInfo struct {
	ID       string
	Date     string
	Notes    string
	Sections []domain.WorkoutSection
}

func MakeWorkoutInfo(workout domain.Workout) WorkoutInfo {
	return WorkoutInfo{
		ID:       string(workout.ID),
		Date:     workout.Date.Format(time.DateOnly),
		Notes:    workout.Notes,
		Sections: workout.Sections,
	}
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

type WorkoutCreateSchema struct {
	Date     string                  `json:"date"`
	Notes    string                  `json:"notes"`
	Sections []domain.WorkoutSection `json:"sections"`
}

func (s *WorkoutCreateSchema) toDomain() (domain.Workout, error) {
	date, err := time.Parse(time.DateOnly, s.Date)
	if err != nil {
		return domain.Workout{}, err
	}

	return domain.Workout{
		Date:     date,
		Notes:    s.Notes,
		Sections: s.Sections,
	}, nil
}

type TrackInfo struct {
	ID   string
	Name string
}

type TrackSchema struct {
	Track   TrackInfo
	IsOwner bool
}
