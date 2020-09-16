package command

import (
	"errors"
	"fmt"
	"github.com/chaseisabelle/goresp"
	"github.com/chaseisabelle/phprom/environment"
	"github.com/chaseisabelle/phprom/funcs"
	"github.com/chaseisabelle/phprom/metrics"
	"time"
)

type Register struct {
	environment *environment.Environment
}

func NewRegister(e *environment.Environment) *Register {
	return &Register{
		environment: e,
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

	args = args[4:]

	var met metrics.Metric

	switch kind[0] {
	case metrics.HISTOGRAM:
		met, err = r.buildHistogram(name, desc, labs, args)
	case metrics.SUMMARY:
		met, err = r.buildSummary(name, desc, labs, args)
	case metrics.GAUGE:
		met, err = r.buildGauge(name, desc, labs, args)
	case metrics.COUNTER:
		met, err = r.buildCounter(name, desc, labs, args)
	default:
		return funcs.RRErrs(fmt.Sprintf("invalid type %+v", kind[0]))
	}

	if err != nil {
		return funcs.ERErrs(err)
	}

	err = met.Register()

	if err != nil {
		return funcs.ERErrs(err)
	}

	return []goresp.Value{goresp.NewSimpleString("ok")}, nil
}

func (r *Register) buildHistogram(name string, desc string, labs []string, args []goresp.Value) (metrics.Metric, error) {
	if len(args) > 1 {
		return nil, errors.New("too many arguments")
	}

	his := metrics.NewHistogram(r.environment)

	his.Name = name
	his.Description = desc
	his.Labels = labs

	if len(args) == 0 {
		return his, nil
	}

	vals, err := args[3].Array()

	if err != nil {
		return his, err
	}

	his.Buckets, err = goresp.Float64s(vals)

	return his, err
}

func (r *Register) buildSummary(name string, desc string, labs []string, args []goresp.Value) (metrics.Metric, error) {
	sum := metrics.NewSummary(r.environment)

	sum.Name = name
	sum.Description = desc
	sum.Labels = labs

	if len(args) == 0 {
		return sum, nil
	}

	var err error

	switch len(args) {
	case 7:
		sum.BufCap, err = args[6].Uint32()

		if err != nil {
			return nil, err
		}

		fallthrough
	case 6:
		sum.Buckets, err = args[5].Uint32()

		if err != nil {
			return nil, err
		}

		fallthrough
	case 5:
		d, err := args[4].Duration(time.Second)

		if err != nil {
			return nil, err
		}

		sum.Age = d

		fallthrough
	case 4:
		vals, err := args[3].Array()

		if err != nil {
			return nil, err
		}

		sum.Objectives = make(map[float64]float64)

		for _, val := range vals {
			a, err := val.Array()

			if err != nil {
				return nil, err
			}

			if len(a) != 2 {
				return nil, errors.New("invalid tuple")
			}

			k, err := a[0].Float64()

			if err != nil {
				return nil, err
			}

			v, err := a[1].Float64()

			if err != nil {
				return nil, err
			}

			sum.Objectives[k] = v
		}
	}

	return sum, nil
}

func (r *Register) buildGauge(name string, desc string, labs []string, args []goresp.Value) (metrics.Metric, error) {
	if len(args) > 0 {
		return nil, errors.New("too many args")
	}

	g := metrics.NewGauge(r.environment)

	g.Name = name
	g.Description = desc
	g.Labels = labs

	return g, nil
}

func (r *Register) buildCounter(name string, desc string, labs []string, args []goresp.Value) (metrics.Metric, error) {
	if len(args) > 0 {
		return nil, errors.New("too many args")
	}

	c := metrics.NewCounter(r.environment)

	c.Name = name
	c.Description = desc
	c.Labels = labs

	return c, nil
}
