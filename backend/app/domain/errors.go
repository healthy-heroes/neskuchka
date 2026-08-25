package domain

import "errors"

var ErrNotFound = errors.New("not found")
var ErrForbidden = errors.New("forbidden")

// ErrLocked means the object's state forbids an operation its caller is
// otherwise allowed to make: the track owner may edit their workouts, but not
// one that has already moved into the past.
var ErrLocked = errors.New("locked")
