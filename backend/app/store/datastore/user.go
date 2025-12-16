package datastore

import (
	"github.com/rs/zerolog/log"

	"github.com/healthy-heroes/neskuchka/backend/app/store"
)

type UserDBStore struct {
	*DataStore
}

func (ds *UserDBStore) Create(user *store.User) (*store.User, error) {
	_, err := ds.Exec(`INSERT INTO user (id, name, email, picture) VALUES (?, ?, ?, ?)`,
		user.ID, user.Name, user.Email, user.Picture)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ds *UserDBStore) Get(id store.UserID) (*store.User, error) {
	user := &store.User{}
	err := ds.DB.Get(user, `SELECT * FROM user WHERE id = ?`, id)

	if err != nil {
		return nil, handleFindError(err)
	}
	return user, nil
}

func (ds *UserDBStore) FindByEmail(email string) (*store.User, error) {
	user := &store.User{}
	err := ds.DB.Get(user, `SELECT * FROM user WHERE email = ?`, email)

	if err != nil {
		return nil, handleFindError(err)
	}
	return user, nil
}

func (ds *UserDBStore) FindOrCreate(email string) (*store.User, error) {
	user, err := ds.FindByEmail(email)
	if err != nil && err != store.ErrNotFound {
		log.Error().Err(err).Msgf("Error while finding user by email %s", email)

		return nil, err
	}

	if user == nil {
		log.Info().Msgf("Creating new user %s", email)

		user, err = ds.User.Create(&store.User{
			ID: store.CreateUserId(),
			// todo: generate new name
			Name:  "New User",
			Email: email,
		})

		if err != nil {
			log.Error().Err(err).Msgf("Error creating user %s", email)

			return nil, err
		}
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
			email TEXT NOT NULL UNIQUE,
			picture TEXT NOT NULL
		)
	`)

	if err != nil {
		log.Error().Err(err).Msg("Failed to create user table")
		return err
	}

	log.Debug().Msg("User table created")
	return nil
}
