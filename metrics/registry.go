package metrics

import (
	"bytes"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"sync"
)

type Registry struct {
	sync.Mutex
	local   *prometheus.Registry
	metrics map[string]Metric
}

var instance *Registry

func registry() *Registry {
	if instance == nil {
		instance = &Registry{
			local:   prometheus.NewRegistry(),
			metrics: make(map[string]Metric),
		}
	}

	return instance
}

func Register(metric Metric) (*prometheus.Collector, error) {
	var collector prometheus.Collector
	var name string

	reg := registry()

	reg.Lock()

	defer reg.Unlock()

	switch metric.Type() {
	case COUNTER:
		counter := metric.(*Counter)
		name = counter.Name
		collector = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    counter.Name,
			Help:    counter.Description,
		}, counter.Labels)
	case GAUGE:
		gauge := metric.(*Gauge)
		name = gauge.Name
		collector = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: gauge.Name,
			Help: gauge.Description,
		}, gauge.Labels)
	case HISTOGRAM:
		histogram := metric.(*Histogram)
		name = histogram.Name
		collector = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    histogram.Name,
			Help:    histogram.Description,
			Buckets: histogram.Buckets,
		}, histogram.Labels)
	case SUMMARY:
		summary := metric.(*Summary)
		name = summary.Name
		collector = prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Name:       summary.Name,
			Help:       summary.Description,
			Objectives: summary.Objectives,
			MaxAge:     summary.Age,
			AgeBuckets: summary.Buckets,
			BufCap:     summary.BufCap,
		}, summary.Labels)
	}

	err := reg.local.Register(collector)

	if err == nil {
		reg.metrics[name] = metric
	}

	return &collector, err
}

func Registered(name string) (Metric, error) {
	reg := registry()

	reg.Lock()

	m, ok := reg.metrics[name]

	reg.Unlock()

	if !ok {
		return nil, fmt.Errorf("unregistered metric: %s", name)
	}

	return m, nil
}

func Metrics() (string, error) {
	gatherers := prometheus.Gatherers{
		registry().local,
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
