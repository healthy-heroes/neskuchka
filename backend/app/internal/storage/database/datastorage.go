package database

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
)

// DataStorage stores common data of the application
type DataStorage struct {
	engine *Engine
	logger zerolog.Logger
}

// NewDataStorage creates a new data storage
func NewDataStorage(engine *Engine, logger zerolog.Logger) *DataStorage {
	return &DataStorage{
		engine: engine,
		logger: logger.With().Str("pkg", "database").Logger(),
	}
}

/**
type dataStorage interface {
	GetUser(context.Context, UserID) (User, error)
	GetUserByEmail(context.Context, Email) (User, error)
	CreateUser(context.Context, User) (User, error)
	UpdateUser(context.Context, User) (User, error)

	GetTrackBySlug(context.Context, TrackSlug) (Track, error)

	GetWorkout(context.Context, WorkoutID) (Workout, error)
	FindWorkouts(context.Context, WorkoutFindCriteria) ([]Workout, error)
	CreateWorkout(context.Context, Workout) (Workout, error)
	UpdateWorkout(context.Context, Workout) (Workout, error)
}
*/

func (ds *DataStorage) GetTrackBySlug(_ context.Context, _ domain.TrackSlug) (domain.Track, error) {
	return domain.Track{}, fmt.Errorf("not implemented")
}

func (ds *DataStorage) GetWorkout(_ context.Context, _ domain.WorkoutID) (domain.Workout, error) {
	return domain.Workout{}, fmt.Errorf("not implemented")
}

func (ds *DataStorage) FindWorkouts(_ context.Context, _ domain.WorkoutFindCriteria) ([]domain.Workout, error) {
	return nil, fmt.Errorf("not implemented")
}
func (ds *DataStorage) CreateWorkout(_ context.Context, _ domain.Workout) (domain.Workout, error) {
	return domain.Workout{}, fmt.Errorf("not implemented")
}
func (ds *DataStorage) UpdateWorkout(_ context.Context, _ domain.Workout) (domain.Workout, error) {
	return domain.Workout{}, fmt.Errorf("not implemented")
}
