package registry

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/chaseisabelle/phprom/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"sync"
	"time"
)

type Registry struct {
	sync.Mutex
	local      *prometheus.Registry
	namespace  string
	histograms map[string]*prometheus.HistogramVec
	counters   map[string]*prometheus.CounterVec
	summaries  map[string]*prometheus.SummaryVec
	gauges     map[string]*prometheus.GaugeVec
}

type Parameters struct {
	Name      string
	Help      string
	Type      byte
	Labels    []string
	Histogram struct {
		Buckets []float64
	}
	Summary struct {
		Objectives map[float64]float64
		Age        time.Duration
		Buckets    uint32
		BufCap     uint32
	}
}

func New(namespace string) (*Registry, error) {
	if namespace == "" {
		return nil, errors.New("empty namespace")
	}

	return &Registry{
		local:      prometheus.NewRegistry(),
		namespace:  namespace,
		histograms: make(map[string]*prometheus.HistogramVec),
		counters:   make(map[string]*prometheus.CounterVec),
		summaries:  make(map[string]*prometheus.SummaryVec),
		gauges:     make(map[string]*prometheus.GaugeVec),
	}, nil
}

func (r *Registry) Register(parameters *Parameters) error {
	var collector prometheus.Collector
	var ok bool

	kind := parameters.Type
	name := parameters.Name

	r.Lock()

	defer r.Unlock()

	switch kind {
	case types.Histogram:
		_, ok = r.histograms[name]
	case types.Counter:
		_, ok = r.counters[name]
	case types.Summary:
		_, ok = r.summaries[name]
	case types.Gauge:
		_, ok = r.gauges[name]
	}

	if ok {
		return nil
	}

	help := parameters.Help
	labels := parameters.Labels

	switch kind {
	case types.Histogram:
		collector = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: r.Namespace(),
			Name:      name,
			Help:      help,
			Buckets:   parameters.Histogram.Buckets,
		}, labels)
	case types.Counter:
		collector = prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: r.Namespace(),
			Name:      name,
			Help:      help,
		}, labels)
	case types.Summary:
		collector = prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Namespace:  r.Namespace(),
			Name:       name,
			Help:       help,
			Objectives: parameters.Summary.Objectives,
			MaxAge:     parameters.Summary.Age,
			AgeBuckets: parameters.Summary.Buckets,
			BufCap:     parameters.Summary.BufCap,
		}, labels)
	case types.Gauge:
		collector = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: r.Namespace(),
			Name:      name,
			Help:      help,
		}, labels)
	}

	err := r.local.Register(collector)

	if err != nil {
		return err
	}

	switch kind {
	case types.Histogram:
		r.histograms[name] = collector.(*prometheus.HistogramVec)
	case types.Counter:
		r.counters[name] = collector.(*prometheus.CounterVec)
	case types.Summary:
		r.summaries[name] = collector.(*prometheus.SummaryVec)
	case types.Gauge:
		r.gauges[name] = collector.(*prometheus.GaugeVec)
	}

	return nil
}

func (r *Registry) Namespace() string {
	return r.namespace
}

func (r *Registry) Histogram(name string) (*prometheus.HistogramVec, error) {
	h, ok := r.histograms[name]

	if !ok {
		return nil, fmt.Errorf("unregistered histogram %s", name)
	}

	return h, nil
}

func (r *Registry) Counter(name string) (*prometheus.CounterVec, error) {
	c, ok := r.counters[name]

	if !ok {
		return nil, fmt.Errorf("unregistered counter %s", name)
	}

	return c, nil
}

func (r *Registry) Summary(name string) (*prometheus.SummaryVec, error) {
	s, ok := r.summaries[name]

	if !ok {
		return nil, fmt.Errorf("unregistered summary %s", name)
	}

	return s, nil
}

func (r *Registry) Gauge(name string) (*prometheus.GaugeVec, error) {
	g, ok := r.gauges[name]

	if !ok {
		return nil, fmt.Errorf("unregistered gauge %s", name)
	}

	return g, nil
}

func (r *Registry) Metrics() (string, error) {
	gatherers := prometheus.Gatherers{
		r.local,
	}

	gathering, err := gatherers.Gather()

	if err != nil {
		return "", err
	}

	out := &bytes.Buffer{}

	for _, family := range gathering {
		_, err := expfmt.MetricFamilyToText(out, family)

		if err != nil {
			break
		}
	}

	return out.String(), err
}
