package cmd

type CommonOptionsCommander interface {
	Execute(args []string) error
}