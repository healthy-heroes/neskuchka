package domain

import (
	"context"

	"github.com/healthy-heroes/neskuchka/backend/app/internal/uuid"
)

type TrackID string
type TrackSlug string

// NewTrackID generates a new track id
func NewTrackID() TrackID {
	return TrackID(uuid.New())
}

// Track is a track aggregate
type Track struct {
	ID          TrackID
	Slug        TrackSlug
	Name        string
	Description string

	OwnerID UserID
}

func (t *Track) IsOwner(userID UserID) bool {
	return userID != "" && t.OwnerID == userID
}

// ApplyUpdate applies an update to a track.
// Only the texts the owner edits are taken; slug and ownership are not theirs
// to change through this path.
func (t *Track) ApplyUpdate(tu Track) {
	t.Name = tu.Name
	t.Description = tu.Description
}

// GetTrack gets a track by id
func (s *Store) GetTrack(ctx context.Context, tid TrackID) (Track, error) {
	return s.storage.GetTrack(ctx, tid)
}

// GetMainTrack gets the main track
func (s *Store) GetMainTrack(ctx context.Context) (Track, error) {
	return s.storage.GetTrackBySlug(ctx, TrackSlug("main"))
}

// UpdateTrack updates the track texts on behalf of its owner
func (s *Store) UpdateTrack(ctx context.Context, uid UserID, tu Track) (Track, error) {
	t, err := s.storage.GetTrack(ctx, tu.ID)
	if err != nil {
		return Track{}, err
	}

	// Permission check
	if !t.IsOwner(uid) {
		return Track{}, ErrForbidden
	}

	t.ApplyUpdate(tu)

	return s.storage.UpdateTrack(ctx, t)
}
