package store

type ExerciseSlug string

type Exercise struct {
	Slug        ExerciseSlug
	Name        string
	Description string
}

type ExerciseStore interface {
	Create(exercise *Exercise) (*Exercise, error)
	Get(slug ExerciseSlug) (*Exercise, error)
	Find(criteria ExerciseFindCriteria) ([]*Exercise, error)
}

type ExerciseFindCriteria struct {
	Limit int
}
