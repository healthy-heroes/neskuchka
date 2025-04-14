package datastore

import (
	"github.com/rs/zerolog/log"

	"github.com/healthy-heroes/neskuchka/backend/app/store"
)

type UserDBStore struct {
	*DataStore
}

func (ds *UserDBStore) Create(user *store.User) (*store.User, error) {
	_, err := ds.Exec(`INSERT INTO user (id, name, login, email) VALUES (?, ?, ?, ?)`,
		user.ID, user.Name, user.Login, user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ds *UserDBStore) Get(id store.UserID) (*store.User, error) {
	user := &store.User{}
	err := ds.DB.Get(user, `SELECT * FROM user WHERE id = ?`, id)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ds *UserDBStore) InitTables() error {
	log.Debug().Msg("Creating user table")

	// Create user table
	_, err := ds.Exec(`
		CREATE TABLE IF NOT EXISTS user (
			id TEXT PRIMARY KEY NOT NULL,
			name TEXT NOT NULL,
			login TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE
		)
	`)

	if err != nil {
		log.Error().Err(err).Msg("Failed to create user table")
		return err
	}

	log.Debug().Msg("User table created")
	return nil
}
