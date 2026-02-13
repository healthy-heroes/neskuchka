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

// GetMainTrack gets the main track
func (s *Store) GetMainTrack(ctx context.Context) (Track, error) {
	return s.dataStorage.GetTrackBySlug(ctx, TrackSlug("main"))
}
