package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"net/http"
	"sync"
)

type Payload struct {
	Metric  string            `json:"metric"`
	Name    string            `json:"name"`
	Help    string            `json:"help"`
	Labels  map[string]string `json:"labels"`
	Buckets []float64         `json:"buckets"`
	Value   float64           `json:"value"`
}

type Registry struct {
	sync.Mutex
	local   *prometheus.Registry
	metrics map[string]string
}

const HISTOGRAM = "histogram"
const COUNTER = "counter"
const SUMMARY = "summary"
const GAUGE = "gauge"

var registry *Registry

func init() {
	registry = &Registry{
		local: prometheus.NewRegistry(),
		metrics: make(map[string]string),
	}
}

func main() {
	host := flag.String("host", ":8080", "host and port to listen on")

	flag.Parse()

	mux := http.NewServeMux()

	mux.HandleFunc("/", handle)

	err := http.ListenAndServe(*host, mux)

	if err != nil {
		panic(err)
	}
}

func handle(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		out, err := metrics()

		if err != nil {
			fail(res, err, http.StatusBadRequest)

			break
		}

		res.WriteHeader(http.StatusOK)

		cnt, err := res.Write([]byte(out))

		if err != nil {
			println(err.Error())

			if cnt != 0 {
				break
			}

			fail(res, err, http.StatusInternalServerError)
		}
	case "POST":
		raw := &Payload{}
		err := json.NewDecoder(req.Body).Decode(raw)

		if err != nil {
			fail(res, err, http.StatusBadRequest)

			break
		}

		switch raw.Metric {
		case HISTOGRAM:
			fail(res, histogram(raw), http.StatusInternalServerError)
		case COUNTER:
			fail(res, counter(raw), http.StatusInternalServerError)
		case SUMMARY:
			fail(res, summary(raw), http.StatusInternalServerError)
		case GAUGE:
			fail(res, gauge(raw), http.StatusInternalServerError)
		default:
			fail(res, fmt.Errorf("invalid metric %s", raw.Metric), http.StatusBadRequest)
		}
	default:
		fail(res, fmt.Errorf("method not allowed %s", req.Method), http.StatusMethodNotAllowed)
	}
}

func fail(res http.ResponseWriter, err error, stat int) {
	if err == nil {
		return
	}

	str := err.Error()

	println(str)

	enc, err := json.Marshal(struct {
		Error string `json:"error"`
	}{
		Error: str,
	})

	if err != nil {
		println(err.Error())

		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	http.Error(res, string(enc), stat)
}

func keys(m map[string]string) []string {
	keys := make([]string, 0, len(m))

	for key := range m {
		keys = append(keys, key)
	}

	return keys
}

func register(col prometheus.Collector, name string, met string) (prometheus.Collector, error) {
	registry.Lock()

	tmp, ok := registry.metrics[name]

	if ok && tmp != met {
		return col, fmt.Errorf("metric %s is already registered as %s not %s", name, tmp, met)
	}

	registry.metrics[name] = met

	registry.Unlock()

	err := registry.local.Register(col)

	if err != nil {
		are, ok := err.(prometheus.AlreadyRegisteredError)

		if !ok {
			return nil, err
		}

		col = are.ExistingCollector
	}

	return col, nil
}

func gauge(raw *Payload) error {
	opts := prometheus.GaugeOpts{
		Name: raw.Name,
		Help: raw.Help,
	}

	g := prometheus.NewGaugeVec(opts, keys(raw.Labels))

	if g == nil {
		return errors.New("failed to init gauge")
	}

	col, err := register(g, raw.Name, raw.Metric)

	if err != nil {
		return err
	}

	g = col.(*prometheus.GaugeVec)

	g.With(raw.Labels).Add(raw.Value)

	return nil
}

func summary(raw *Payload) error {
	opts := prometheus.SummaryOpts{
		Name: raw.Name,
		Help: raw.Help,
	}

	sum := prometheus.NewSummaryVec(opts, keys(raw.Labels))

	if sum == nil {
		return errors.New("failed to init summary")
	}

	col, err := register(sum, raw.Name, raw.Metric)

	if err != nil {
		return err
	}

	sum = col.(*prometheus.SummaryVec)

	sum.With(raw.Labels).Observe(raw.Value)

	return nil
}

func counter(raw *Payload) error {
	opts := prometheus.CounterOpts{
		Name: raw.Name,
		Help: raw.Help,
	}

	cnt := prometheus.NewCounterVec(opts, keys(raw.Labels))

	if cnt == nil {
		return errors.New("failed to init counter")
	}

	col, err := register(cnt, raw.Name, raw.Metric)

	if err != nil {
		return err
	}

	cnt = col.(*prometheus.CounterVec)

	cnt.With(raw.Labels).Add(raw.Value)

	return nil
}

func histogram(raw *Payload) error {
	opts := prometheus.HistogramOpts{
		Name: raw.Name,
		Help: raw.Help,
	}

	if raw.Buckets != nil && len(raw.Buckets) > 0 {
		opts.Buckets = raw.Buckets
	}

	his := prometheus.NewHistogramVec(opts, keys(raw.Labels))

	if his == nil {
		return errors.New("failed to init histogram")
	}

	col, err := register(his, raw.Name, raw.Metric)

	if err != nil {
		return err
	}

	his = col.(*prometheus.HistogramVec)

	his.With(raw.Labels).Observe(raw.Value)

	return nil
}

func metrics() (string, error) {
	mfs, err := (prometheus.Gatherers{
		registry.local,
	}).Gather()

	if err != nil {
		return "", err
	}

	out := &bytes.Buffer{}

	for _, mf := range mfs {
		_, err = expfmt.MetricFamilyToText(out, mf)

		if err != nil {
			return "", err
		}
	}

	return out.String(), nil
}
