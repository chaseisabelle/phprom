package v1

import (
	"bytes"
	"context"
	"fmt"
	phprom_v1 "github.com/chaseisabelle/phprom/pkg/v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
)

type PHProm struct{}

var registry *prometheus.Registry
var counters map[string]*prometheus.CounterVec
var histograms map[string]*prometheus.HistogramVec
var summaries map[string]*prometheus.SummaryVec
var gauges map[string]*prometheus.GaugeVec

func init() {
	registry = prometheus.NewRegistry()
	counters = make(map[string]*prometheus.CounterVec)
	histograms = make(map[string]*prometheus.HistogramVec)
	summaries = make(map[string]*prometheus.SummaryVec)
	gauges = make(map[string]*prometheus.GaugeVec)
}

func New() (*PHProm, error) {
	return &PHProm{}, nil
}

func (p *PHProm) Get(ctx context.Context, req *phprom_v1.GetRequest) (*phprom_v1.GetResponse, error) {
	mfs, err := prometheus.Gatherers{
		registry,
	}.Gather()

	if err != nil {
		return nil, err
	}

	out := &bytes.Buffer{}

	for _, fam := range mfs {
		_, err := expfmt.MetricFamilyToText(out, fam)

		if err != nil {
			return nil, err
		}
	}

	return &phprom_v1.GetResponse{
		Metrics: out.String(),
	}, nil
}

func (p *PHProm) RegisterCounter(ctx context.Context, req *phprom_v1.RegisterCounterRequest) (*phprom_v1.RegisterResponse, error) {
	col := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: req.Namespace,
		Name:      req.Name,
		Help:      req.Description,
	}, req.Labels)

	res, err := register(col)

	if err == nil && !res.Registered {
		counters[key(req.Namespace, req.Name)] = col
	}

	return res, err
}

func (p *PHProm) RegisterHistogram(ctx context.Context, req *phprom_v1.RegisterHistogramRequest) (*phprom_v1.RegisterResponse, error) {
	bux := make([]float64, len(req.Buckets))

	for i, b := range req.Buckets {
		bux[i] = float64(b)
	}

	col := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: req.Namespace,
		Name:      req.Name,
		Help:      req.Description,
		Buckets:   bux,
	}, req.Labels)

	res, err := register(col)

	if err == nil && !res.Registered {
		histograms[key(req.Namespace, req.Name)] = col
	}

	return res, err
}

func (p *PHProm) RegisterSummary(ctx context.Context, req *phprom_v1.RegisterSummaryRequest) (*phprom_v1.RegisterResponse, error) {
	col := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: req.Namespace,
		Name:      req.Name,
		Help:      req.Description,
	}, req.Labels)

	res, err := register(col)

	if err == nil && !res.Registered {
		summaries[key(req.Namespace, req.Name)] = col
	}

	return res, err
}

func (p *PHProm) RegisterGauge(ctx context.Context, req *phprom_v1.RegisterGaugeRequest) (*phprom_v1.RegisterResponse, error) {
	col := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: req.Namespace,
		Name:      req.Name,
		Help:      req.Description,
	}, req.Labels)

	res, err := register(col)

	if err == nil && !res.Registered {
		gauges[key(req.Namespace, req.Name)] = col
	}

	return res, err
}

func (p *PHProm) RecordCounter(ctx context.Context, req *phprom_v1.RecordCounterRequest) (*phprom_v1.RecordResponse, error) {
	col, ok := counters[key(req.Namespace, req.Name)]

	if !ok {
		return nil, fmt.Errorf("no counter registered as %s", req.Name)
	}

	vec, err := col.GetMetricWith(req.Labels)

	if err != nil {
		return nil, err
	}

	vec.Add(float64(req.Value))

	return &phprom_v1.RecordResponse{}, nil
}

func (p *PHProm) RecordHistogram(ctx context.Context, req *phprom_v1.RecordHistogramRequest) (*phprom_v1.RecordResponse, error) {
	col, ok := histograms[key(req.Namespace, req.Name)]

	if !ok {
		return nil, fmt.Errorf("no histogram registered as %s", req.Name)
	}

	vec, err := col.GetMetricWith(req.Labels)

	if err != nil {
		return nil, err
	}

	vec.Observe(float64(req.Value))

	return &phprom_v1.RecordResponse{}, nil
}

func (p *PHProm) RecordSummary(ctx context.Context, req *phprom_v1.RecordSummaryRequest) (*phprom_v1.RecordResponse, error) {
	col, ok := summaries[key(req.Namespace, req.Name)]

	if !ok {
		return nil, fmt.Errorf("no summary registered as %s", req.Name)
	}

	vec, err := col.GetMetricWith(req.Labels)

	if err != nil {
		return nil, err
	}

	vec.Observe(float64(req.Value))

	return &phprom_v1.RecordResponse{}, nil
}

func (p *PHProm) RecordGauge(ctx context.Context, req *phprom_v1.RecordGaugeRequest) (*phprom_v1.RecordResponse, error) {
	col, ok := gauges[key(req.Namespace, req.Name)]

	if !ok {
		return nil, fmt.Errorf("no gauge registered as %s", req.Name)
	}

	vec, err := col.GetMetricWith(req.Labels)

	if err != nil {
		return nil, err
	}

	vec.Add(float64(req.Value))

	return &phprom_v1.RecordResponse{}, nil
}

func key(ns string, n string) string {
	return fmt.Sprintf("%s_%s", ns, n)
}

func register(c prometheus.Collector) (*phprom_v1.RegisterResponse, error) {
	err := registry.Register(c)

	_, ok := err.(prometheus.AlreadyRegisteredError)

	if ok {
		err = nil
	}

	return &phprom_v1.RegisterResponse{
		Registered: ok,
	}, err
}
