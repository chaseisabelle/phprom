package metrics

import (
	"github.com/chaseisabelle/phprom/environment"
	"github.com/prometheus/client_golang/prometheus"
)

type Histogram struct {
	Name        string
	Description string
	Labels      []string
	Buckets     []float64
	environment *environment.Environment
	internal    *prometheus.HistogramVec
}

func NewHistogram(e *environment.Environment) *Histogram {
	return &Histogram{
		environment: e,
	}
}

func (h *Histogram) Type() byte {
	return 'h'
}

func (h *Histogram) Register() error {
	err := canRegisterAs(h.Name, h.Type())

	if err != nil {
		return err
	}

	col, err := Register(h)

	if err != nil {
		return err
	}

	h.internal = (*col).(*prometheus.HistogramVec)

	return nil
}

func (h *Histogram) Record(value float64, labels map[string]string) error {
	h.internal.With(labels).Observe(value)

	return nil
}
