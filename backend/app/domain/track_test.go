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
			Storage: &StorageStub{
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

func TestTrack_IsOwner(t *testing.T) {
	t.Run("should return true if user is the owner", func(t *testing.T) {
		track := Track{
			OwnerID: UserID("user-1"),
		}
		assert.True(t, track.IsOwner(UserID("user-1")))
	})
	t.Run("should return false if user is not the owner", func(t *testing.T) {
		track := Track{
			OwnerID: UserID("user-1"),
		}
		assert.False(t, track.IsOwner(UserID("user-2")))
		assert.False(t, track.IsOwner(UserID("")))
	})
}

func TestTrack_ApplyUpdate(t *testing.T) {
	t.Run("should take the texts and nothing else", func(t *testing.T) {
		track := createTrack()
		track.Name = "Main"
		track.Description = "The main track"

		track.ApplyUpdate(Track{
			ID:          NewTrackID(),
			Slug:        TrackSlug("hijacked"),
			Name:        "Renamed",
			Description: "Rewritten",
			OwnerID:     NewUserID(),
		})

		assert.Equal(t, "Renamed", track.Name)
		assert.Equal(t, "Rewritten", track.Description)

		assert.Equal(t, TrackSlug("main-test"), track.Slug)
		assert.NotEmpty(t, track.OwnerID)
	})
}

func TestUpdateTrack(t *testing.T) {
	setup := func(track Track) *Store {
		return NewStore(Opts{
			Storage: &StorageStub{
				GetTrackFunc: func(ctx context.Context, tid TrackID) (Track, error) {
					return track, nil
				},
				UpdateTrackFunc: func(ctx context.Context, t Track) (Track, error) {
					return t, nil
				},
			},
		})
	}

	t.Run("should update the texts for the owner", func(t *testing.T) {
		track := createTrack()

		updated, err := setup(track).UpdateTrack(context.Background(), track.OwnerID, Track{
			ID:          track.ID,
			Name:        "Renamed",
			Description: "Rewritten",
		})

		assert.NoError(t, err)
		assert.Equal(t, "Renamed", updated.Name)
		assert.Equal(t, "Rewritten", updated.Description)
		assert.Equal(t, track.OwnerID, updated.OwnerID)
	})

	t.Run("should refuse a stranger", func(t *testing.T) {
		track := createTrack()

		_, err := setup(track).UpdateTrack(context.Background(), NewUserID(), Track{ID: track.ID})

		assert.ErrorIs(t, err, ErrForbidden)
	})

	t.Run("should refuse an anonymous caller", func(t *testing.T) {
		track := createTrack()

		_, err := setup(track).UpdateTrack(context.Background(), UserID(""), Track{ID: track.ID})

		assert.ErrorIs(t, err, ErrForbidden)
	})
}
