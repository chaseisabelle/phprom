package command

import (
	"errors"
	"github.com/chaseisabelle/goresp"
	"github.com/chaseisabelle/phprom/environment"
	"github.com/chaseisabelle/phprom/funcs"
	"github.com/chaseisabelle/phprom/registry"
	"github.com/chaseisabelle/phprom/types"
)

type Register struct {
	env *environment.Environment
}

func NewRegister(e *environment.Environment) *Register {
	return &Register{
		env: e,
	}
}

func (r *Register) Execute(args ...goresp.Value) ([]goresp.Value, error) {
	var err error

	switch len(args) {
	case 0:
		err = errors.New("metric name required")
	case 1:
		err = errors.New("metric description required")
	case 2:
		err = errors.New("metric type required")
	case 3:
		err = errors.New("metric labels required")
	}

	if err != nil {
		return funcs.ERErrs(err)
	}

	name, err := args[0].String()

	if err != nil {
		return funcs.ERErrs(err)
	}

	desc, err := args[1].String()

	if err != nil {
		return funcs.ERErrs(err)
	}

	kind, err := args[2].Bytes()

	if err == nil && len(kind) != 1 {
		err = errors.New("malformed metric type")
	}

	if err != nil {
		return funcs.ERErrs(err)
	}

	vals, err := args[3].Array()

	if err != nil {
		return funcs.ERErrs(err)
	}

	labs := make([]string, len(vals))

	for i, val := range vals {
		lab, err := val.String()

		if err != nil {
			return funcs.ERErrs(err)
		}

		labs[i] = lab
	}

	pars := &registry.Parameters{}

	pars.Name = name
	pars.Help = desc
	pars.Type = kind[0]
	pars.Labels = labs

	switch kind[0] {
	case types.Histogram:
		if len(args) > 4 {
			return funcs.RRErrs("too many arguments")
		}

		if len(args) == 4 {
			bux, err := args[3]
		}

		pars.Histogram.Buckets
	}

	reg, err := r.env.Registry()

	if err != nil {
		return funcs.ERErrs(err)
	}

	reg.Register(pars)
}
