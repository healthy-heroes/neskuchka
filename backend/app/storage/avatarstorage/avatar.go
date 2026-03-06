package avatarstorage

import (
	"context"
	"time"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
	"github.com/healthy-heroes/neskuchka/backend/app/storage"
)

type avatarRow struct {
	UserID string `db:"user_id"`

	MimeType string `db:"mime_type"`
	Data     []byte `db:"data"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func makeAvatar(id domain.UserID, a domain.Avatar) avatarRow {
	return avatarRow{
		UserID:   string(id),
		MimeType: a.MimeType,
		Data:     a.Data,

		UpdatedAt: time.Now(),
	}
}

func (a avatarRow) toDomain() domain.Avatar {
	return domain.Avatar{
		MimeType: a.MimeType,
		Data:     a.Data,
	}
}

func (s *Storage) Get(ctx context.Context, id domain.UserID) (domain.Avatar, error) {
	avatar := avatarRow{}

	err := s.engine.GetContext(ctx, &avatar, "SELECT * FROM avatar WHERE user_id = ?", id)
	if err != nil {
		return domain.Avatar{}, storage.HandleSqlError(err)
	}
	return avatar.toDomain(), nil
}

func (s *Storage) Save(ctx context.Context, id domain.UserID, avatar domain.Avatar) error {
	row := makeAvatar(id, avatar)

	_, err := s.engine.ExecContext(ctx, `
		INSERT INTO avatar (user_id, mime_type, data, updated_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(user_id) DO UPDATE SET
			mime_type = excluded.mime_type,
			data = excluded.data,
			updated_at = excluded.updated_at
	`, row.UserID, row.MimeType, row.Data, row.UpdatedAt)
	if err != nil {
		return storage.HandleSqlError(err)
	}

	return nil
}

func (s *Storage) Delete(ctx context.Context, id domain.UserID) error {
	_, err := s.engine.ExecContext(ctx, "DELETE FROM avatar WHERE user_id = ?", id)
	return storage.HandleSqlError(err)
}

func (s *Storage) Exists(ctx context.Context, id domain.UserID) (bool, error) {
	var exists bool
	err := s.engine.GetContext(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM avatar WHERE user_id = ?)", id)
	if err != nil {
		return false, storage.HandleSqlError(err)
	}
	return exists, nil
}
