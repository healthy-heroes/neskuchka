package domain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TrackStoreStub struct {
	GetBySlugFunc func(context.Context, TrackSlug) (Track, error)
}

func (s TrackStoreStub) GetBySlug(ctx context.Context, slug TrackSlug) (Track, error) {
	return s.GetBySlugFunc(ctx, slug)
}

func TestNewTrackID(t *testing.T) {
	t.Run("should generate a new track id", func(t *testing.T) {
		trackID := NewTrackID()
		assert.NotEmpty(t, trackID)
	})
}

func TestGetMainTrack(t *testing.T) {
	t.Run("should return main track", func(t *testing.T) {
		service := NewService(Opts{
			TrackStore: &TrackStoreStub{
				GetBySlugFunc: func(ctx context.Context, slug TrackSlug) (Track, error) {
					return Track{
						ID:   TrackID("1"),
						Slug: slug,
					}, nil
				},
			},
		})
		track, err := service.GetMainTrack(context.Background())

		assert.Nil(t, err)
		assert.Equal(t, track.Slug, TrackSlug("main"))
	})
}
