package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
)

func getUserDb(t *testing.T, engine *Engine, id string) userDb {
	user := userDb{}
	err := engine.Get(&user, "SELECT * FROM user WHERE id = ?", id)
	require.NoError(t, err)

	return user
}

func Test_User_Create(t *testing.T) {
	ds := setupTestDataStorage(t, setupTestSqliteDB(t))
	defer ds.engine.Close()

	user := domain.User{
		ID:    domain.NewUserID(),
		Name:  "Test User",
		Email: "test@example.com",
	}

	created, err := ds.CreateUser(context.Background(), user)
	require.NoError(t, err)
	assert.Equal(t, user, created)

	getByID, err := ds.GetUser(context.Background(), user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user, getByID)

	getByEmail, err := ds.GetUserByEmail(context.Background(), user.Email)
	assert.NoError(t, err)
	assert.Equal(t, user, getByEmail)

	// checks system fields
	userDb := getUserDb(t, ds.engine, string(created.ID))
	assert.NotZero(t, userDb.CreatedAt)
	assert.NotZero(t, userDb.UpdatedAt)
}

func Test_User_Update(t *testing.T) {
	ds := setupTestDataStorage(t, setupTestSqliteDB(t))
	defer ds.engine.Close()

	existingUser := domain.User{
		ID:    domain.NewUserID(),
		Name:  "Test User",
		Email: "test@example.com",
	}
	_, err := ds.CreateUser(context.Background(), existingUser)
	require.NoError(t, err)

	createdUser := getUserDb(t, ds.engine, string(existingUser.ID))

	updateUser := domain.User{
		ID:    existingUser.ID,
		Name:  "Test User Updated",
		Email: "wrong@example.com",
	}
	u, err := ds.UpdateUser(context.Background(), updateUser)
	require.NoError(t, err)
	assert.NotEqual(t, updateUser, u)
	assert.Equal(t, existingUser.ID, u.ID)
	assert.Equal(t, existingUser.Email, u.Email)

	// no created wrong user
	_, err = ds.GetUserByEmail(context.Background(), updateUser.Email)
	assert.Equal(t, domain.ErrNotFound, err)

	// checks system fields
	updatedUser := getUserDb(t, ds.engine, string(u.ID))
	assert.Equal(t, createdUser.CreatedAt, updatedUser.CreatedAt)
	assert.Greater(t, updatedUser.UpdatedAt, createdUser.UpdatedAt)
}
