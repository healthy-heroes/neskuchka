package datastore

import (
	"testing"

	"github.com/healthy-heroes/neskuchka/backend/app/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTrackTestDB(t *testing.T) *DataStore {
	ds := setupTestDB(t)

	// Initialize user table first because track has a foreign key to user
	err := ds.User.InitTables()
	require.NoError(t, err)

	err = ds.Track.InitTables()
	require.NoError(t, err)

	return ds
}

func createTestUser(t *testing.T, ds *DataStore) *store.User {
	userID := store.CreateUserId()
	user := &store.User{
		ID:      userID,
		Name:    "Test User",
		Email:   "test@example.com",
		Picture: "test.png",
	}
	_, err := ds.User.Create(user)
	require.NoError(t, err)
	return user
}

func TestTrackDBStore_Create(t *testing.T) {
	ds := setupTrackTestDB(t)
	defer ds.Close()

	user := createTestUser(t, ds)

	trackID := store.CreateTrackId()
	track := &store.Track{
		ID:      trackID,
		Name:    "Test Track",
		OwnerID: user.ID,
	}

	// Test creating a new track
	created, err := ds.Track.Create(track)
	require.NoError(t, err)
	assert.Equal(t, track, created)
}

func TestTrackDBStore_Get(t *testing.T) {
	ds := setupTrackTestDB(t)
	defer ds.Close()

	user := createTestUser(t, ds)

	trackID := store.CreateTrackId()
	track := &store.Track{
		ID:      trackID,
		Name:    "Test Track",
		OwnerID: user.ID,
	}

	_, err := ds.Track.Create(track)
	require.NoError(t, err)

	// Test getting an existing track
	found, err := ds.Track.Get(track.ID)
	require.NoError(t, err)
	assert.Equal(t, track, found)

	// Test getting a non-existent track
	nonExistentID := store.CreateTrackId()
	_, err = ds.Track.Get(nonExistentID)
	assert.Error(t, err, "Should error when track not found")
}

func TestTrackDBStore_GetMainTrack(t *testing.T) {
	ds := setupTrackTestDB(t)
	defer ds.Close()

	user := createTestUser(t, ds)

	// Create two tracks
	track1 := &store.Track{
		ID:      store.CreateTrackId(),
		Name:    "Track 1",
		OwnerID: user.ID,
	}

	track2 := &store.Track{
		ID:      store.CreateTrackId(),
		Name:    "Track 2",
		OwnerID: user.ID,
	}

	_, err := ds.Track.Create(track1)
	require.NoError(t, err)

	_, err = ds.Track.Create(track2)
	require.NoError(t, err)

	// Get main track should return the first one
	mainTrack, err := ds.Track.GetMainTrack()
	require.NoError(t, err)
	assert.NotNil(t, mainTrack)
}
