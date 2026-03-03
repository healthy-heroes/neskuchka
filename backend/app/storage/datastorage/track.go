package datastorage

import (
	"context"
	"time"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
)

type trackRow struct {
	ID          string
	Slug        string
	Name        string
	Description string

	OwnerId string `db:"owner_id"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func makeTrack(t domain.Track) trackRow {
	return trackRow{
		ID:          string(t.ID),
		Slug:        string(t.Slug),
		Name:        t.Name,
		Description: t.Description,
		OwnerId:     string(t.OwnerID),

		UpdatedAt: time.Now(),
	}
}

func (t trackRow) toDomain() domain.Track {
	return domain.Track{
		ID:          domain.TrackID(t.ID),
		Slug:        domain.TrackSlug(t.Slug),
		Name:        t.Name,
		Description: t.Description,
		OwnerID:     domain.UserID(t.OwnerId),
	}
}

func (s *Storage) GetTrack(ctx context.Context, id domain.TrackID) (domain.Track, error) {
	t := trackRow{}
	err := s.engine.Get(&t, "SELECT * FROM track WHERE id = ?", id)

	return t.toDomain(), handleSqlError(err)
}

func (s *Storage) GetTrackBySlug(ctx context.Context, slug domain.TrackSlug) (domain.Track, error) {
	t := trackRow{}
	err := s.engine.Get(&t, "SELECT * FROM track WHERE slug = ?", slug)

	return t.toDomain(), handleSqlError(err)
}

func (s *Storage) CreateTrack(ctx context.Context, track domain.Track) (domain.Track, error) {
	t := makeTrack(track)

	_, err := s.engine.Exec(
		"INSERT INTO track(id, slug, name, description, owner_id) VALUES(?,?,?,?,?)",
		t.ID, t.Slug, t.Name, t.Description, t.OwnerId,
	)

	if err != nil {
		return domain.Track{}, handleSqlError(err)
	}

	return s.GetTrack(ctx, track.ID)
}
