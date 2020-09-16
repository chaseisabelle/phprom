package metrics

import (
	"github.com/chaseisabelle/phprom/environment"
	"github.com/prometheus/client_golang/prometheus"
)

type Gauge struct {
	Name        string
	Description string
	Labels      []string
	environment *environment.Environment
	internal    *prometheus.GaugeVec
}

func NewGauge(e *environment.Environment) *Gauge {
	return &Gauge{
		environment: e,
	}
}

func (g *Gauge) Type() byte {
	return 'g'
}

func (g *Gauge) Register() error {
	err := canRegisterAs(g.Name, g.Type())

	if err != nil {
		return err
	}

	col, err := Register(g)

	if err != nil {
		return err
	}

	g.internal = (*col).(*prometheus.GaugeVec)

	return nil
}

func (g *Gauge) Record(value float64, labels map[string]string) error {
	g.internal.With(labels).Add(value)

	return nil
}
