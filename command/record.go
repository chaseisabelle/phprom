package command

import (
	"github.com/chaseisabelle/goresp"
	"github.com/chaseisabelle/phprom/environment"
	"github.com/chaseisabelle/phprom/funcs"
	"github.com/chaseisabelle/phprom/metrics"
)

type Record struct {
	environment *environment.Environment
}

func NewRecord(e *environment.Environment) *Record {
	return &Record{
		environment: e,
	}
}

func (r *Record) Execute(args ...goresp.Value) ([]goresp.Value, error) {
	if len(args) < 2 {
		return funcs.RRErrs("not enough args")
	}

	if len(args) > 3 {
		return funcs.RRErrs("too many args")
	}

	str, err := args[0].String()

	if err != nil {
		return funcs.ERErrs(err)
	}

	met, err := metrics.Registered(str)

	if err != nil {
		return funcs.ERErrs(err)
	}

	flt, err := args[1].Float64()

	if err != nil {
		return funcs.ERErrs(err)
	}

	var labs map[string]string

	if len(args) == 3 {
		vals, err := args[2].Array()

		if err != nil {
			return funcs.ERErrs(err)
		}

		labs, err = mapLabels(vals)

		if err != nil {
			return funcs.ERErrs(err)
		}
	}

	err = met.Record(flt, labs)

	if err != nil {
		return funcs.ERErrs(err)
	}

	return []goresp.Value{goresp.NewSimpleString("ok")}, nil
}
