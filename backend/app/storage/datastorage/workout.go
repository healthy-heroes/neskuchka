package datastorage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
	"github.com/healthy-heroes/neskuchka/backend/app/storage"
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

func (s *Storage) GetWorkout(ctx context.Context, wr domain.WorkoutRef) (domain.Workout, error) {
	workout := workoutRow{}
	err := s.engine.GetContext(ctx, &workout, "SELECT * FROM workout WHERE track_id = ? AND id = ?", wr.TrackID, wr.WorkoutID)
	if err != nil {
		return domain.Workout{}, storage.HandleSqlError(err)
	}

	return workout.toDomain()
}

// workoutOrder is the one ordering of a track, and the keyset predicate below
// has to match it exactly. Ties on date break by id: ids are UUIDv7, so they
// sort by creation, which is what created_at was reaching for — except
// created_at is only second-resolution and leaves same-second rows unordered.
const workoutOrder = " ORDER BY date DESC, id DESC"

func (s *Storage) FindWorkouts(ctx context.Context, tid domain.TrackID, criteria domain.WorkoutFindCriteria, now time.Time) ([]domain.Workout, error) {
	query := "SELECT * FROM workout WHERE track_id = ?"
	args := []any{tid}

	if criteria.PublishedOnly {
		query += " AND date <= ?"
		args = append(args, now.Format(time.DateOnly))
	}

	if !criteria.After.IsZero() {
		query += " AND (date < ? OR (date = ? AND id < ?))"
		cursorDate := criteria.After.Date.Format(time.DateOnly)
		args = append(args, cursorDate, cursorDate, criteria.After.ID)
	}

	query += workoutOrder + " LIMIT ?"
	args = append(args, criteria.Limit)

	workouts := []workoutRow{}
	err := s.engine.SelectContext(ctx, &workouts, query, args...)
	if err != nil {
		return nil, storage.HandleSqlError(err)
	}

	return rowsToDomain(workouts)
}

// CountWorkouts counts the track in one pass: everything in it, and the part
// still ahead of now.
func (s *Storage) CountWorkouts(ctx context.Context, tid domain.TrackID, now time.Time) (int, int, error) {
	counts := struct {
		Total   int
		Planned int
	}{}

	err := s.engine.GetContext(ctx, &counts,
		"SELECT COUNT(*) AS total, COALESCE(SUM(date > ?), 0) AS planned FROM workout WHERE track_id = ?",
		now.Format(time.DateOnly), tid,
	)
	if err != nil {
		return 0, 0, storage.HandleSqlError(err)
	}

	return counts.Total, counts.Planned, nil
}


func (s *Storage) CreateWorkout(ctx context.Context, workout domain.Workout) (domain.Workout, error) {
	w, err := makeWorkout(workout)
	if err != nil {
		return domain.Workout{}, err
	}

	_, err = s.engine.ExecContext(ctx,
		"INSERT INTO workout(id, track_id, date, sections, notes, schema_version) VALUES(?,?,?,?,?,?)",
		w.ID, w.TrackID, w.Date, w.Sections, w.Notes, workoutSchemaVersion,
	)
	if err != nil {
		return domain.Workout{}, storage.HandleSqlError(err)
	}

	return s.GetWorkout(ctx, domain.WorkoutRef{TrackID: workout.TrackID, WorkoutID: workout.ID})
}

func (s *Storage) DeleteWorkout(ctx context.Context, wr domain.WorkoutRef) error {
	result, err := s.engine.ExecContext(ctx,
		"DELETE FROM workout WHERE track_id = ? AND id = ?",
		wr.TrackID, wr.WorkoutID,
	)
	if err != nil {
		return storage.HandleSqlError(err)
	}

	// DELETE is happy to remove nothing, so the caller has to be able to tell a
	// deleted row from one that was never there.
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (s *Storage) UpdateWorkout(ctx context.Context, workout domain.Workout) (domain.Workout, error) {
	w, err := makeWorkout(workout)
	if err != nil {
		return domain.Workout{}, err
	}

	_, err = s.engine.ExecContext(ctx,
		"UPDATE workout SET date = ?, sections = ?, notes = ?, updated_at = ? WHERE track_id = ? AND id = ?",
		w.Date, w.Sections, w.Notes, w.UpdatedAt,
		w.TrackID, w.ID,
	)
	if err != nil {
		return domain.Workout{}, storage.HandleSqlError(err)
	}

	return s.GetWorkout(ctx, domain.WorkoutRef{TrackID: workout.TrackID, WorkoutID: workout.ID})
}
