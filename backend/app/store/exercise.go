package store

type ExerciseSlug string

type Exercise struct {
	Slug        ExerciseSlug
	Name        string
	Description string
}

type ExerciseStore interface {
	Store

	Create(exercise *Exercise) (*Exercise, error)
	Get(slug ExerciseSlug) (*Exercise, error)
	Find(criteria *ExerciseFindCriteria) ([]*Exercise, error)
}

type ExerciseFindCriteria struct {
	Slugs []ExerciseSlug
	Limit int
}
