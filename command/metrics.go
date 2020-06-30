package command

import (
	"errors"
	"github.com/chaseisabelle/goresp"
	"github.com/chaseisabelle/phprom/environment"
)

type Metrics struct {
	env *environment.Environment
}

func NewMetrics(e *environment.Environment) *Metrics {
	return &Metrics{
		env: e,
	}
}

func (m *Metrics) Execute(args ...goresp.Value) ([]goresp.Value, error) {
	if len(args) > 0 {
		err := errors.New("metrics command does not accept arguments")

		return []goresp.Value{goresp.NewError(err)}, err
	}

	reg, err := m.env.Registry()

	if err != nil {
		return []goresp.Value{goresp.NewError(err)}, err
	}

	met, err := reg.Metrics()

	if err != nil {
		return []goresp.Value{goresp.NewError(err)}, err
	}

	return []goresp.Value{goresp.NewBulkString(met)}, nil
}
