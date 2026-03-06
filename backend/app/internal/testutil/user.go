package testutil

import (
	"fmt"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
)

func CreateUser() domain.User {
	id := domain.NewUserID()
	email := domain.Email(fmt.Sprintf("user-%s@example.com", id))

	return domain.User{
		ID:    id,
		Name:  fmt.Sprintf("User #%s", id),
		Email: email,
	}
}
