package store

import "github.com/google/uuid"

type TrackID string

type Track struct {
	ID      TrackID
	Name    string
	OwnerID UserID
}

func (t *Track) IsOwner(userID UserID) bool {
	return userID != "" && t.OwnerID == userID
}

func CreateTrackId() TrackID {
	return TrackID(uuid.New().String())
}

type TrackStore interface {
	Store

	Create(track *Track) (*Track, error)
	Get(id TrackID) (*Track, error)
	GetMainTrack() (*Track, error)
}
