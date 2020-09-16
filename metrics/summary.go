package metrics

import (
	"github.com/chaseisabelle/phprom/environment"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type Summary struct {
	Name        string
	Description string
	Labels      []string
	Objectives  map[float64]float64
	Age         time.Duration
	Buckets     uint32
	BufCap      uint32
	environment *environment.Environment
	internal    *prometheus.SummaryVec
}

func NewSummary(e *environment.Environment) *Summary {
	return &Summary{
		environment: e,
	}
}

func (s *Summary) Type() byte {
	return 's'
}

func (s *Summary) Register() error {
	err := canRegisterAs(s.Name, s.Type())

	if err != nil {
		return err
	}

	col, err := Register(s)

	if err != nil {
		return err
	}

	s.internal = (*col).(*prometheus.SummaryVec)

	return nil
}

func (s *Summary) Record(value float64, labels map[string]string) error {
	s.internal.With(labels).Observe(value)

	return nil
}
