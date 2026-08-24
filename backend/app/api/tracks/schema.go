package tracks

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
)

// WorkoutInfo is a workout info for reading and creating workouts
type WorkoutInfo struct {
	ID       string
	TrackID  string
	Date     string
	Notes    string
	Sections []domain.WorkoutSection

	// IsPublished and CanEdit are the workout's state, not its data: whether
	// participants can see it yet and whether the owner may still change it.
	// They ride along so that the client shows the same answer the domain would
	// give, instead of reimplementing the rule next to the buttons.
	IsPublished bool
	CanEdit     bool
}

func MakeWorkoutInfo(workout domain.Workout) WorkoutInfo {
	now := time.Now()

	return WorkoutInfo{
		ID:       string(workout.ID),
		TrackID:  string(workout.TrackID),
		Date:     workout.Date.Format(time.DateOnly),
		Notes:    workout.Notes,
		Sections: makeSections(workout.Sections),

		IsPublished: workout.IsPublished(now),
		CanEdit:     workout.IsEditable(now),
	}
}

// makeSections holds the API contract that every list is an array, never null.
//
// An exercise without a prescription is ordinary — that is most of a warm-up.
// A section without exercises is not a state the product allows, and the editor
// refuses to save one; the guard stays for rows written before that was settled.
// A nil slice marshals to null, and clients iterate these lists without checking.
func makeSections(sections []domain.WorkoutSection) []domain.WorkoutSection {
	result := make([]domain.WorkoutSection, 0, len(sections))

	for _, section := range sections {
		exercises := make([]domain.WorkoutExercise, 0, len(section.Exercises))
		for _, exercise := range section.Exercises {
			if exercise.Prescription == nil {
				exercise.Prescription = []string{}
			}
			exercises = append(exercises, exercise)
		}

		section.Exercises = exercises
		result = append(result, section)
	}

	return result
}

func (w *WorkoutInfo) toDomain() (domain.Workout, error) {
	date, err := time.Parse(time.DateOnly, w.Date)
	if err != nil {
		return domain.Workout{}, err
	}

	return domain.Workout{
		ID:       domain.WorkoutID(w.ID),
		TrackID:  domain.TrackID(w.TrackID),
		Date:     date,
		Notes:    w.Notes,
		Sections: w.Sections,
	}, nil
}

func MakeWorkoutInfos(workouts []domain.Workout) []WorkoutInfo {
	workoutInfos := make([]WorkoutInfo, 0, len(workouts))
	for _, workout := range workouts {
		workoutInfos = append(workoutInfos, MakeWorkoutInfo(workout))
	}
	return workoutInfos
}

type WorkoutSchema struct {
	Workout WorkoutInfo
}

type WorkoutsSchema struct {
	Workouts []WorkoutInfo
}

// WorkoutsPageSchema is a page of the track for its owner: the rows asked for,
// plus the counters the list header shows above them.
type WorkoutsPageSchema struct {
	Workouts []WorkoutInfo

	Total   int
	Planned int
}

type TrackInfo struct {
	ID          string
	Name        string
	Description string

	Author AuthorInfo
}

// AuthorInfo is the track owner shown next to the track description.
// No avatar yet: the user service only emits an avatar URL when the file
// actually exists, and wiring the avatar store in here is not worth it.
type AuthorInfo struct {
	ID   string
	Name string
}

func MakeTrackInfo(track domain.Track, author domain.User) TrackInfo {
	return TrackInfo{
		ID:          string(track.ID),
		Name:        track.Name,
		Description: track.Description,
		Author: AuthorInfo{
			ID:   string(author.ID),
			Name: author.Name,
		},
	}
}

type TrackSchema struct {
	Track   TrackInfo
	IsOwner bool
}

// UpdateTrackSchema is the payload of the track texts the owner edits.
// Only these two: neither the slug nor the ownership is editable from outside.
type UpdateTrackSchema struct {
	Name        string
	Description string
}

func (s UpdateTrackSchema) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&s.Description, validation.Length(0, 2000)),
	)
}

func (s UpdateTrackSchema) toDomain(id domain.TrackID) domain.Track {
	return domain.Track{
		ID:          id,
		Name:        s.Name,
		Description: s.Description,
	}
}
