package cmd

import "fmt"

type CommonOptionsCommander interface {
	SetCommon(opts *CommonOptions)
	Execute(args []string) error
}

type CommonOptions struct {
	Revision string
}

func (c *CommonOptions) SetCommon(opts *CommonOptions) {
	c.Revision = opts.Revision
}

func (opts *CommonOptions) Execute(args []string) error {
	return fmt.Errorf("execute method does not implemented")
}
