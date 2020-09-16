package metrics

import (
	"github.com/chaseisabelle/phprom/environment"
	"github.com/prometheus/client_golang/prometheus"
)

type Counter struct {
	Name        string
	Description string
	Labels      []string
	environment *environment.Environment
	internal    *prometheus.CounterVec
}

func NewCounter(e *environment.Environment) *Counter {
	return &Counter{
		environment: e,
	}
}

func (c *Counter) Type() byte {
	return 'c'
}

func (c *Counter) Register() error {
	err := canRegisterAs(c.Name, c.Type())

	if err != nil {
		return err
	}

	col, err := Register(c)

	if err != nil {
		return err
	}

	c.internal = (*col).(*prometheus.CounterVec)

	return nil
}

func (c *Counter) Record(value float64, labels map[string]string) error {
	c.internal.With(labels).Add(value)

	return nil
}
