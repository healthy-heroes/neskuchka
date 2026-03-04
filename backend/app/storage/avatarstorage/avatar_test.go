package avatarstorage

import (
	"testing"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_AvatarStorage_SaveAndGet(t *testing.T) {
	storage := setupTestStorage(t)

	userID := domain.NewUserID()
	avatar := domain.Avatar{
		MimeType: "image/png",
		Data:     []byte("test"),
	}

	err := storage.Save(t.Context(), userID, avatar)
	require.NoError(t, err)

	avatarByID, err := storage.Get(t.Context(), userID)
	require.NoError(t, err)
	assert.Equal(t, avatar, avatarByID)

	_, err = storage.Get(t.Context(), domain.NewUserID())
	assert.ErrorIs(t, err, domain.ErrNotFound)
}
