package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
)

func trackFromDB(t *testing.T, engine *Engine, id string) trackRow {
	track := trackRow{}
	err := engine.Get(&track, "SELECT * FROM track WHERE id = ?", id)
	require.NoError(t, err)

	return track
}

func Test_Track_Create(t *testing.T) {
	ds := setupTestDataStorage(t)
	defer ds.engine.Close()

	newTrack := domain.Track{
		ID:          domain.NewTrackID(),
		Slug:        domain.TrackSlug("testmain"),
		Name:        "Test track",
		Description: "Its track created for tests",
		OwnerID:     domain.UserID("user-1"),
	}

	createdTrack, err := ds.CreateTrack(context.Background(), newTrack)
	require.NoError(t, err)
	assert.Equal(t, newTrack, createdTrack)

	trackByID, err := ds.GetTrack(context.Background(), newTrack.ID)
	require.NoError(t, err)
	assert.Equal(t, newTrack, trackByID)

	trackBySlug, err := ds.GetTrackBySlug(context.Background(), newTrack.Slug)
	require.NoError(t, err)
	assert.Equal(t, newTrack, trackBySlug)

	// checks system fields
	createdRow := trackFromDB(t, ds.engine, string(newTrack.ID))
	assert.NotZero(t, createdRow.CreatedAt)
	assert.NotZero(t, createdRow.UpdatedAt)
}
