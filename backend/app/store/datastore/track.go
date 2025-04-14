package datastore

import (
	"github.com/rs/zerolog/log"

	"github.com/healthy-heroes/neskuchka/backend/app/store"
)

type TrackDBStore struct {
	*DataStore
}

type TrackDB struct {
	ID      store.TrackID `db:"id"`
	Name    string        `db:"name"`
	OwnerID store.UserID  `db:"owner_id"`
}

func (t *TrackDB) ToStore() *store.Track {
	return &store.Track{
		ID:      t.ID,
		Name:    t.Name,
		OwnerID: t.OwnerID,
	}
}

func (ds *TrackDBStore) Create(track *store.Track) (*store.Track, error) {
	_, err := ds.Exec(`INSERT INTO track (id, name, owner_id) VALUES (?, ?, ?)`,
		track.ID, track.Name, track.OwnerID)
	if err != nil {
		return nil, err
	}
	return track, nil
}

func (ds *TrackDBStore) Get(id store.TrackID) (*store.Track, error) {
	track := &TrackDB{}
	err := ds.DB.Get(track, `SELECT * FROM track WHERE id = ?`, id)

	if err != nil {
		return nil, err
	}
	return track.ToStore(), nil
}

func (ds *TrackDBStore) GetMainTrack() (*store.Track, error) {
	// We'll use the first track as the main track for now
	// In a real application, you might have a "main" flag or similar
	track := &TrackDB{}
	err := ds.DB.Get(track, `SELECT * FROM track ORDER BY id LIMIT 1`)

	if err != nil {
		return nil, err
	}
	return track.ToStore(), nil
}

func (ds *TrackDBStore) InitTables() error {
	log.Debug().Msg("Creating track table")

	// Create track table
	_, err := ds.Exec(`
		CREATE TABLE IF NOT EXISTS track (
			id TEXT PRIMARY KEY NOT NULL,
			name TEXT NOT NULL,
			owner_id TEXT NOT NULL
		)
	`)

	if err != nil {
		log.Error().Err(err).Msg("Failed to create track table")
		return err
	}

	log.Debug().Msg("Track table created")
	return nil
}
