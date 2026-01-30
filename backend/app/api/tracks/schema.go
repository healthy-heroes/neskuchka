package tracks

import "github.com/healthy-heroes/neskuchka/backend/app/store"

type WorkoutSchema struct {
	Workout *store.Workout
}

type WorkoutsSchema struct {
	Workouts []*store.Workout
}

type TrackInfo struct {
	ID   string
	Name string
}

type TrackSchema struct {
	Track   TrackInfo
	IsOwner bool
}
