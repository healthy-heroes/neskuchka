package devcmd

type DevCommand struct {
	SeedCmd  SeedCommand  `command:"seed" description:"Seed the database with initial data"`
	TokenCmd TokenCommand `command:"token" description:"Generate a token for the user"`
}
