package datastorage

import (
	"context"
	"time"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
	"github.com/healthy-heroes/neskuchka/backend/app/storage"
)

type userRow struct {
	ID    string
	Email string
	Name  string

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func makeUser(u domain.User) userRow {
	return userRow{
		ID:    string(u.ID),
		Email: string(u.Email),
		Name:  u.Name,

		UpdatedAt: time.Now(),
	}
}

func (u userRow) toDomain() domain.User {
	return domain.User{
		ID:    domain.UserID(u.ID),
		Email: domain.Email(u.Email),
		Name:  u.Name,
	}
}

func (s *Storage) GetUser(ctx context.Context, id domain.UserID) (domain.User, error) {
	u := userRow{}
	err := s.engine.GetContext(ctx, &u, "SELECT * FROM user WHERE id = ?", id)

	return u.toDomain(), storage.HandleSqlError(err)
}

func (s *Storage) GetUserByEmail(ctx context.Context, email domain.Email) (domain.User, error) {
	u := userRow{}
	err := s.engine.GetContext(ctx, &u, "SELECT * FROM user WHERE email = ?", email)

	return u.toDomain(), storage.HandleSqlError(err)
}

func (s *Storage) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	u := makeUser(user)

	_, err := s.engine.ExecContext(ctx, "INSERT INTO user(id, email, name) VALUES(?,?,?)",
		u.ID, u.Email, u.Name)

	if err != nil {
		return domain.User{}, storage.HandleSqlError(err)
	}

	return s.GetUser(ctx, user.ID)
}

func (s *Storage) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	u := makeUser(user)

	_, err := s.engine.ExecContext(ctx, "UPDATE user SET name = ?, updated_at = ? WHERE id = ?", u.Name, u.UpdatedAt, u.ID)

	if err != nil {
		return domain.User{}, storage.HandleSqlError(err)
	}

	return s.GetUser(ctx, user.ID)
}
