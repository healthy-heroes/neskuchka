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

// TrackStore is a interface for track storage
type TrackStore interface {
	GetBySlug(context.Context, TrackSlug) (Track, error)
}

func (s *Service) GetMainTrack(ctx context.Context) (Track, error) {
	return s.trackStore.GetBySlug(ctx, TrackSlug("main"))
}
