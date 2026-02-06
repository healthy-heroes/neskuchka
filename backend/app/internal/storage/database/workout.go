package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
)

const workoutSchemaVersion = 1

type workoutRow struct {
	ID       string
	TrackID  string `db:"track_id"`
	Date     string
	Sections []byte
	Notes    string

	SchemaVersion int       `db:"schema_version"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

// makeWorkout converts a domain.Workout to a workoutRow
func makeWorkout(w domain.Workout) (workoutRow, error) {
	sectionsJSON, err := json.Marshal(w.Sections)
	if err != nil {
		return workoutRow{}, fmt.Errorf("failed to marshal sections: %w", err)
	}

	return workoutRow{
		ID:       string(w.ID),
		TrackID:  string(w.TrackID),
		Date:     w.Date.Format(time.DateOnly),
		Sections: sectionsJSON,
		Notes:    w.Notes,

		UpdatedAt: time.Now(),
	}, nil
}

// toDomain converts a workoutRow to a domain.Workout
func (w *workoutRow) toDomain() (domain.Workout, error) {
	sections := []domain.WorkoutSection{}
	err := json.Unmarshal(w.Sections, &sections)
	if err != nil {
		return domain.Workout{}, fmt.Errorf("failed to unmarshal sections: %w", err)
	}

	date, err := time.Parse(time.DateOnly, w.Date)
	if err != nil {
		return domain.Workout{}, fmt.Errorf("failed to parse date: %w", err)
	}

	return domain.Workout{
		ID:       domain.WorkoutID(w.ID),
		TrackID:  domain.TrackID(w.TrackID),
		Date:     date,
		Sections: sections,
		Notes:    w.Notes,
	}, nil
}

// rowsToDomain converts a []workoutRow to []domain.Workout, handling errors (fail fast)
func rowsToDomain(rows []workoutRow) ([]domain.Workout, error) {
	workouts := make([]domain.Workout, 0, len(rows))
	for _, w := range rows {
		workout, err := w.toDomain()
		if err != nil {
			return nil, err
		}

		workouts = append(workouts, workout)
	}
	return workouts, nil
}

// GetWorkout returns a workout by id
func (ds *DataStorage) GetWorkout(ctx context.Context, id domain.WorkoutID) (domain.Workout, error) {
	workout := workoutRow{}
	err := ds.engine.Get(&workout, "SELECT * FROM workout WHERE id = ?", id)
	if err != nil {
		return domain.Workout{}, handleSqlError(err)
	}

	return workout.toDomain()
}

// FindWorkouts returns workouts filtered by criteria
func (ds *DataStorage) FindWorkouts(ctx context.Context, criteria domain.WorkoutFindCriteria) ([]domain.Workout, error) {
	workouts := []workoutRow{}
	err := ds.engine.Select(&workouts,
		"SELECT * FROM workout WHERE track_id = ? ORDER BY date DESC LIMIT ?",
		criteria.TrackID, criteria.Limit,
	)
	if err != nil {
		return nil, handleSqlError(err)
	}

	return rowsToDomain(workouts)
}

// CreateWorkout creates a new workout and returns it
func (ds *DataStorage) CreateWorkout(ctx context.Context, workout domain.Workout) (domain.Workout, error) {
	w, err := makeWorkout(workout)
	if err != nil {
		return domain.Workout{}, err
	}

	_, err = ds.engine.Exec(
		"INSERT INTO workout(id, track_id, date, sections, notes, schema_version) VALUES(?,?,?,?,?,?)",
		w.ID, w.TrackID, w.Date, w.Sections, w.Notes, workoutSchemaVersion,
	)
	if err != nil {
		return domain.Workout{}, handleSqlError(err)
	}

	return ds.GetWorkout(ctx, workout.ID)
}

// UpdateWorkout updates a workout and returns it
func (ds *DataStorage) UpdateWorkout(ctx context.Context, workout domain.Workout) (domain.Workout, error) {
	w, err := makeWorkout(workout)
	if err != nil {
		return domain.Workout{}, err
	}

	_, err = ds.engine.Exec(
		"UPDATE workout SET date = ?, sections = ?, notes = ?, updated_at = ? WHERE track_id = ? AND id = ?",
		w.Date, w.Sections, w.Notes, w.UpdatedAt,
		w.TrackID, w.ID,
	)
	if err != nil {
		return domain.Workout{}, handleSqlError(err)
	}

	return ds.GetWorkout(ctx, workout.ID)
}
