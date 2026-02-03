package domain

type TrackID string
type TrackSlug string

type Track struct {
	ID          TrackID
	Slug        TrackSlug
	Name        string
	Description string

	OwnerID UserID
}
