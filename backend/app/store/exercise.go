package store

type ExerciseSlug string

type Exercise struct {
	Slug ExerciseSlug `json:"slug"`
	Name string       `json:"name"`
}

type ExerciseStore interface {
	Create(exercise *Exercise) error
	Get(slug ExerciseSlug) (*Exercise, error)
	Find(criteria ExerciseGetCriteria) ([]*Exercise, error)
}

type ExerciseGetCriteria struct {
	Limit   int
	AfterID int
}
