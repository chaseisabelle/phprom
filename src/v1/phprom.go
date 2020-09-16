package v1

import (
	"bytes"
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
)

type PHProm struct{}

var counters map[string]*prometheus.CounterVec
var histograms map[string]*prometheus.HistogramVec
var summaries map[string]*prometheus.SummaryVec
var gauges map[string]*prometheus.GaugeVec

func init() {
	counters = make(map[string]*prometheus.CounterVec)
	histograms = make(map[string]*prometheus.HistogramVec)
	summaries = make(map[string]*prometheus.SummaryVec)
	gauges = make(map[string]*prometheus.GaugeVec)
}

func New() (*PHProm, error) {
	return &PHProm{}, nil
}

func (p *PHProm) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	gat, err := prometheus.Gatherers{
		prometheus.DefaultGatherer,
	}.Gather()

	if err != nil {
		return nil, err
	}

	out := &bytes.Buffer{}

	for _, fam := range gat {
		_, err := expfmt.MetricFamilyToText(out, fam)

		if err != nil {
			return nil, err
		}
	}

	return &GetResponse{
		Metrics: out.String(),
	}, nil
}

func (p *PHProm) RegisterCounter(ctx context.Context, req *RegisterCounterRequest) (*RegisterResponse, error) {
	col := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: req.Name,
		Help: req.Description,
	}, req.Labels)

	res, err := register(col)

	if err == nil && !res.Registered {
		counters[req.Name] = col
	}

	return res, err
}

func (p *PHProm) RegisterHistogram(ctx context.Context, req *RegisterHistogramRequest) (*RegisterResponse, error) {
	bux := make([]float64, len(req.Buckets))

	for i, b := range req.Buckets {
		bux[i] = float64(b)
	}

	col := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    req.Name,
		Help:    req.Description,
		Buckets: bux,
	}, req.Labels)

	res, err := register(col)

	if err == nil && !res.Registered {
		histograms[req.Name] = col
	}

	return res, err
}

func (p *PHProm) RegisterSummary(ctx context.Context, req *RegisterSummaryRequest) (*RegisterResponse, error) {
	col := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: req.Name,
		Help: req.Description,
	}, req.Labels)

	res, err := register(col)

	if err == nil && !res.Registered {
		summaries[req.Name] = col
	}

	return res, err
}

func (p *PHProm) RegisterGauge(ctx context.Context, req *RegisterGaugeRequest) (*RegisterResponse, error) {
	col := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: req.Name,
		Help: req.Description,
	}, req.Labels)

	res, err := register(col)

	if err == nil && !res.Registered {
		gauges[req.Name] = col
	}

	return res, err
}

func (p *PHProm) RecordCounter(ctx context.Context, req *RecordCounterRequest) (*RecordResponse, error) {
	col, ok := counters[req.Name]

	if !ok {
		return nil, fmt.Errorf("no counter registered as %s", req.Name)
	}

	col.With(req.Labels).Add(float64(req.Value))

	return &RecordResponse{
		Api: "1",
	}, nil
}

func (p *PHProm) RecordHistogram(ctx context.Context, req *RecordHistogramRequest) (*RecordResponse, error) {
	col, ok := histograms[req.Name]

	if !ok {
		return nil, fmt.Errorf("no histogram registered as %s", req.Name)
	}

	col.With(req.Labels).Observe(float64(req.Value))

	return &RecordResponse{
		Api: "1",
	}, nil
}

func (p *PHProm) RecordSummary(ctx context.Context, req *RecordSummaryRequest) (*RecordResponse, error) {
	col, ok := summaries[req.Name]

	if !ok {
		return nil, fmt.Errorf("no summary registered as %s", req.Name)
	}

	col.With(req.Labels).Observe(float64(req.Value))

	return &RecordResponse{
		Api: "1",
	}, nil
}

func (p *PHProm) RecordGauge(ctx context.Context, req *RecordGaugeRequest) (*RecordResponse, error) {
	col, ok := gauges[req.Name]

	if !ok {
		return nil, fmt.Errorf("no gauge registered as %s", req.Name)
	}

	col.With(req.Labels).Add(float64(req.Value))

	return &RecordResponse{
		Api: "1",
	}, nil
}

func register(c prometheus.Collector) (*RegisterResponse, error) {
	err := prometheus.DefaultRegisterer.Register(c)

	_, ok := err.(prometheus.AlreadyRegisteredError)

	if ok {
		err = nil
	}

	return &RegisterResponse{
		Api:        "1",
		Registered: ok,
	}, err
}
