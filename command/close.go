package command

import (
	"errors"
	"github.com/chaseisabelle/goresp"
	"github.com/chaseisabelle/phprom/environment"
	"io"
)

type Close struct {
	env *environment.Environment
}

func NewClose(e *environment.Environment) *Close {
	return &Close{
		env: e,
	}
}

func (c *Close) Execute(args ...goresp.Value) ([]goresp.Value, error) {
	err := io.EOF

	if len(args) > 0 {
		err = errors.New("close command does not accept arguments")
	}

	return []goresp.Value{goresp.NewError(err)}, io.EOF
}

