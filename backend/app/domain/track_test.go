package domain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTrackID(t *testing.T) {
	t.Run("should generate a new track id", func(t *testing.T) {
		trackID := NewTrackID()
		assert.NotEmpty(t, trackID)
	})
}

func TestGetMainTrack(t *testing.T) {
	t.Run("should return main track", func(t *testing.T) {
		service := NewStore(Opts{
			DataStorage: &DataStorageStub{
				GetTrackBySlugFunc: func(ctx context.Context, slug TrackSlug) (Track, error) {
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
