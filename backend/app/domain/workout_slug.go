package domain

import (
	"fmt"
	"strings"
	"time"
)

// WorkoutSlug is how a workout is addressed in a URL: "2026-09-21-3f4a5b", the
// date and the tail of the id.
//
// The date is the part people read; the tail tells two workouts of one day
// apart. Nothing is stored: the slug is derived from the workout, so moving the
// date moves the slug and the old link stops resolving. Unique within a track,
// not across tracks — a track is the namespace, the same as for exercises.
type WorkoutSlug string

// slugTailLen is how much of the id goes into the slug. Ids are UUIDv7, whose
// head is a timestamp and near-identical for anything created in the same
// weeks; the tail is the random part.
//
// Six digits rather than four because the odds have to be counted over every
// same-day pair a track ever accumulates, not over one pair: at two workouts a
// day, four digits put a collision somewhere in a decade at around five
// percent, six digits at a fiftieth of that. Nothing enforces uniqueness — on
// a collision the lower id answers and the other workout has no working link.
const slugTailLen = 6

// WorkoutSlugRef addresses a workout the way a link does: the track it belongs
// to, and what its slug decodes into. The track is part of the address rather
// than a filter on it — a slug is only unique inside one, the same as WorkoutRef
// is only meaningful against its track.
type WorkoutSlugRef struct {
	TrackID TrackID
	Date    time.Time
	IDTail  string
}

func (w *Workout) Slug() WorkoutSlug {
	id := string(w.ID)
	tail := id[len(id)-min(len(id), slugTailLen):]

	return WorkoutSlug(w.Date.Format(time.DateOnly) + "-" + tail)
}

// ParseWorkoutSlug decodes a slug taken from a URL into the workout it names
// in the given track.
//
// A malformed slug is ErrNotFound rather than a validation error: it names
// nothing, the same as a well-formed slug nobody has.
func ParseWorkoutSlug(tid TrackID, slug WorkoutSlug) (WorkoutSlugRef, error) {
	notFound := func() (WorkoutSlugRef, error) {
		return WorkoutSlugRef{}, fmt.Errorf("%w: bad slug %q", ErrNotFound, slug)
	}

	s := string(slug)
	if len(s) != len(time.DateOnly)+1+slugTailLen {
		return notFound()
	}

	date, err := time.Parse(time.DateOnly, s[:len(time.DateOnly)])
	if err != nil {
		return notFound()
	}

	tail := s[len(time.DateOnly)+1:]
	if s[len(time.DateOnly)] != '-' || !isLowerHex(tail) {
		return notFound()
	}

	return WorkoutSlugRef{TrackID: tid, Date: date, IDTail: tail}, nil
}

// isLowerHex matches the alphabet of a uuid string, so that a tail compares to
// the id byte for byte with no normalizing on either side.
func isLowerHex(s string) bool {
	return strings.Trim(s, "0123456789abcdef") == ""
}
