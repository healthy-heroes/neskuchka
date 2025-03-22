package datastore

import (
	"testing"

	"github.com/healthy-heroes/neskuchka/backend/app/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupUserTestDB(t *testing.T) *DataStore {
	ds := setupTestDB(t)
	err := ds.User.InitTables()
	require.NoError(t, err)
	return ds
}

func TestUserDBStore_Create(t *testing.T) {
	ds := setupUserTestDB(t)
	defer ds.Close()

	userID := store.CreateUserId()
	user := &store.User{
		ID:    userID,
		Name:  "Test User",
		Login: "testuser",
		Email: "test@example.com",
	}

	// Test creating a new user
	created, err := ds.User.Create(user)
	require.NoError(t, err)
	assert.Equal(t, user, created)

	// Test duplicate creation should fail
	_, err = ds.User.Create(user)
	assert.Error(t, err, "Should fail on duplicate ID/login/email")
}

func TestUserDBStore_Get(t *testing.T) {
	ds := setupUserTestDB(t)
	defer ds.Close()

	// Create test data
	userID := store.CreateUserId()
	user := &store.User{
		ID:    userID,
		Name:  "Test User",
		Login: "testuser",
		Email: "test@example.com",
	}
	_, err := ds.User.Create(user)
	require.NoError(t, err)

	// Test getting an existing user
	found, err := ds.User.Get(user.ID)
	require.NoError(t, err)
	assert.Equal(t, user, found)

	// Test getting a non-existent user
	nonExistentID := store.CreateUserId()
	_, err = ds.User.Get(nonExistentID)
	assert.Error(t, err, "Should error when user not found")
}
