package api_user

// UserSchema is the schema for data about logged user
type UserSchema struct {
	ID     string
	Name   string
	Avatar string `json:",omitempty"`
}
