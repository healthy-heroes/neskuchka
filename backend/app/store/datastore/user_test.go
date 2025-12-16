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
		ID:      userID,
		Name:    "Test User",
		Email:   "test@example.com",
		Picture: "test.png",
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
		ID:      userID,
		Name:    "Test User",
		Picture: "test.png",
		Email:   "test@example.com",
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

func TestUserDBStore_GetByEmail(t *testing.T) {
	ds := setupUserTestDB(t)
	defer ds.Close()

	// Create test data
	userID := store.CreateUserId()
	user := &store.User{
		ID:      userID,
		Name:    "Test User",
		Picture: "test.png",
		Email:   "test@example.com",
	}
	_, err := ds.User.Create(user)
	require.NoError(t, err)

	_, err = ds.User.Create(&store.User{
		ID:    store.CreateUserId(),
		Name:  "Test User 2",
		Email: "test2@example.com",
	})
	require.NoError(t, err)

	// Test getting an existing user
	found, err := ds.User.FindByEmail("test@example.com")
	require.NoError(t, err)
	assert.Equal(t, user, found)

	// Test getting a non-existent user
	nonExistentID := store.CreateUserId()
	_, err = ds.User.Get(nonExistentID)
	assert.Error(t, err, "Should error when user not found")
}

func TestUserDBStore_FindOrCreate(t *testing.T) {
	ds := setupUserTestDB(t)
	defer ds.Close()

	// Create a user to test finding
	existingUser := &store.User{
		ID:      store.CreateUserId(),
		Name:    "Existing User",
		Picture: "existing.png",
		Email:   "existing@example.com",
	}
	_, err := ds.User.Create(existingUser)
	require.NoError(t, err)

	// Should find the existing user (not create a new one)
	found, err := ds.User.FindOrCreate(existingUser.Email)
	require.NoError(t, err)
	assert.Equal(t, existingUser, found)

	// Should create a new user if not found
	newEmail := "newuser@example.com"
	created, err := ds.User.FindOrCreate(newEmail)
	require.NoError(t, err)
	assert.Equal(t, newEmail, created.Email)
	assert.NotEmpty(t, created.ID)
	assert.NotEqual(t, existingUser.ID, created.ID)

	// Should find the newly created user if called again
	foundAgain, err := ds.User.FindOrCreate(newEmail)
	require.NoError(t, err)
	assert.Equal(t, created, foundAgain)
}
