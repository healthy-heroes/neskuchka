package domain

import "errors"

var ErrNotFound = errors.New("not found")
var ErrForbidden = errors.New("forbidden")

// ErrLocked is the state of an object refusing an operation its caller is
// otherwise allowed to make: the track owner may edit their workouts, but not
// a workout that has already moved into the past.
var ErrLocked = errors.New("locked")
