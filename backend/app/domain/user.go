package domain

type UserID string
type Email string

type User struct {
	ID    UserID
	Name  string
	Email Email
}
