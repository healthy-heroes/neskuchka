package uuid

import "github.com/google/uuid"

// New generates a new UUID v7
func New() string {
	uuid, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}

	return uuid.String()
}
