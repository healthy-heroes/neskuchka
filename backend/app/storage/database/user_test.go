package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
)

func userFromDB(t *testing.T, engine *Engine, id string) userRow {
	user := userRow{}
	err := engine.Get(&user, "SELECT * FROM user WHERE id = ?", id)
	require.NoError(t, err)

	return user
}

func Test_User_Create(t *testing.T) {
	ds := setupTestDataStorage(t)

	newUser := domain.User{
		ID:    domain.NewUserID(),
		Name:  "Test User",
		Email: "test@example.com",
	}

	createdUser, err := ds.CreateUser(t.Context(), newUser)
	require.NoError(t, err)
	assert.Equal(t, newUser, createdUser)

	userByID, err := ds.GetUser(t.Context(), newUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, newUser, userByID)

	userByEmail, err := ds.GetUserByEmail(t.Context(), newUser.Email)
	assert.NoError(t, err)
	assert.Equal(t, newUser, userByEmail)

	// checks system fields
	createdRow := userFromDB(t, ds.engine, string(createdUser.ID))
	assert.NotZero(t, createdRow.CreatedAt)
	assert.NotZero(t, createdRow.UpdatedAt)
}

func Test_User_Update(t *testing.T) {
	ds := setupTestDataStorage(t)
	defer ds.engine.Close()

	existingUser := domain.User{
		ID:    domain.NewUserID(),
		Name:  "Test User",
		Email: "test@example.com",
	}
	_, err := ds.CreateUser(t.Context(), existingUser)
	require.NoError(t, err)

	createdUserRow := userFromDB(t, ds.engine, string(existingUser.ID))

	updateUser := domain.User{
		ID:    existingUser.ID,
		Name:  "Test User Updated",
		Email: "wrong@example.com",
	}
	u, err := ds.UpdateUser(t.Context(), updateUser)
	require.NoError(t, err)
	assert.NotEqual(t, updateUser, u)
	assert.Equal(t, existingUser.ID, u.ID)
	assert.Equal(t, existingUser.Email, u.Email)

	// no created wrong user
	_, err = ds.GetUserByEmail(t.Context(), updateUser.Email)
	assert.Equal(t, domain.ErrNotFound, err)

	// checks system fields
	updatedUserRow := userFromDB(t, ds.engine, string(u.ID))
	assert.Equal(t, createdUserRow.CreatedAt, updatedUserRow.CreatedAt)
	assert.Greater(t, updatedUserRow.UpdatedAt, createdUserRow.UpdatedAt)
}

func Test_User_NotFound(t *testing.T) {
	ds := setupTestDataStorage(t)

	_, err := ds.GetUser(t.Context(), domain.UserID("non-existent-id"))
	assert.ErrorIs(t, err, domain.ErrNotFound)

	_, err = ds.GetUserByEmail(t.Context(), domain.Email("non-existent-email"))
	assert.ErrorIs(t, err, domain.ErrNotFound)
}
