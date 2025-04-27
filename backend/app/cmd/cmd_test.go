package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCommonOptions struct {
	Name string

	CommonOptions
}

func (opts *TestCommonOptions) Execute(args []string) error {
	return nil
}

func TestCommonOptions_Default(t *testing.T) {
	opts := &TestCommonOptions{Name: "default"}
	opts.SetCommon(&CommonOptions{Revision: "test"})

	assert.Equal(t, "test", opts.Revision, "revision should be set")
	assert.Nil(t, opts.Execute([]string{}), "execute should be overridden")
}
