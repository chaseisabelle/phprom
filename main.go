package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"github.com/tidwall/resp"
	"io"
	"net"
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

	lis, err := net.Listen("tcp4", *host)

	if err != nil {
		panic(err.Error())
	}

	defer func() {
		err = lis.Close()

		if err != nil {
			printerr(err)
		}
	}()

	for {
		con, err := lis.Accept()

		if err != nil {
			printerr(err)

			continue
		}

		go func() {
			err := handle(con)

			if err != nil {
				printerr(err)
			}
		}()
	}
}

func printerr(err error) {
	println(err.Error())
}

func handle(con net.Conn) error {
	defer func() {
		err := con.Close()

		if err != nil {
			printerr(err)
		}
	}()

	rdr := resp.NewReader(bufio.NewReader(con))

	for {
		val, _, err := rdr.ReadValue()
		println(fmt.Sprintf("%+v", val))

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return failure(con, err)
		}

		err = respond(con, "test")

		if err != nil {
			return failure(con, err)
		}
	}
}

func respond(con net.Conn, res string) error {
	buf := bytes.Buffer{}
	wri := resp.NewWriter(&buf)
	err := wri.WriteString(res)

	if err != nil {
		return err
	}

	_, err = con.Write(buf.Bytes())

	return err
}

func failure(con net.Conn, err error) error {
	printerr(err)

	return respond(con, err.Error())
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
