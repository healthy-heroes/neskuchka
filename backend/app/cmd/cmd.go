package cmd

import "fmt"

// CommonOptionsCommander is a base interface for commands
type CommonOptionsCommander interface {
	SetCommon(opts *CommonOptions)
	Execute(args []string) error
}

// CommonOptions contains common options for all commands
type CommonOptions struct {
	Revision string
}

// SetCommon sets the common options
func (c *CommonOptions) SetCommon(opts *CommonOptions) {
	c.Revision = opts.Revision
}

// Execute is the main method for executing commands
func (opts *CommonOptions) Execute(args []string) error {
	return fmt.Errorf("execute method does not implemented")
}
