package database

import (
	"context"
	"time"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
)

type userDb struct {
	ID    string
	Email string
	Name  string

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func fromDomain(u domain.User) userDb {
	return userDb{
		ID:    string(u.ID),
		Email: string(u.Email),
		Name:  u.Name,

		UpdatedAt: time.Now(),
	}
}

func (u userDb) toDomain() domain.User {
	return domain.User{
		ID:    domain.UserID(u.ID),
		Email: domain.Email(u.Email),
		Name:  u.Name,
	}
}

// GetUser returns a user by id
func (ds *DataStorage) GetUser(ctx context.Context, id domain.UserID) (domain.User, error) {
	u := userDb{}
	err := ds.engine.Get(&u, "SELECT * FROM user WHERE id = ?", id)

	return u.toDomain(), handleSqlError(err)
}

// GetUserByEmail returns a user by email
func (ds *DataStorage) GetUserByEmail(ctx context.Context, email domain.Email) (domain.User, error) {
	u := userDb{}
	err := ds.engine.Get(&u, "SELECT * FROM user WHERE email = ?", email)

	return u.toDomain(), handleSqlError(err)
}

// CreateUser creates a new user
func (ds *DataStorage) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	u := fromDomain(user)

	_, err := ds.engine.Exec("INSERT INTO user(id, email, name) VALUES(?,?,?)",
		u.ID, u.Email, u.Name)

	if err != nil {
		return domain.User{}, handleSqlError(err)
	}

	return ds.GetUser(ctx, user.ID)
}

// UpdateUser updates user's mutable fields
func (ds *DataStorage) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	u := fromDomain(user)

	_, err := ds.engine.Exec("UPDATE user SET name = ?, updated_at = ? WHERE id = ?", u.Name, u.UpdatedAt, u.ID)

	if err != nil {
		return domain.User{}, handleSqlError(err)
	}

	return ds.GetUser(ctx, user.ID)
}
