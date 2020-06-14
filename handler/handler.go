package handler

import (
	"errors"
	"fmt"
	"github.com/chaseisabelle/phprom/registry"
	"github.com/chaseisabelle/phprom/types"
	"github.com/tidwall/resp"
	"time"
)

type Handler struct {
	local *registry.Registry
}

const (
	Metrics   byte = 'M'
	Register  byte = 'R'
	Histogram      = types.Histogram
	Counter        = types.Counter
	Summary        = types.Summary
	Gauge          = types.Gauge
)

func New(local *registry.Registry) *Handler {
	return &Handler{
		local: local,
	}
}

func (h *Handler) Handle(command byte, arguments []resp.Value) resp.Value {
	var err error

	switch command {
	case Metrics:
		out, err := h.metrics()

		if err == nil {
			return resp.StringValue(out)
		}
	case Register:
		err = h.register(arguments)
	case Histogram:
		fallthrough
	case Counter:
		fallthrough
	case Summary:
		fallthrough
	case Gauge:
		err = h.record(command, arguments)
	default:
		err = fmt.Errorf("invalid main command %b", command)
	}

	if err != nil {
		return resp.ErrorValue(err)
	}

	return resp.NullValue()
}

func (h *Handler) register(arguments []resp.Value) error {
	count := len(arguments)

	if count < 1 {
		return errors.New("name required")
	}

	name := arguments[0].String()

	if count < 2 {
		return errors.New("type required")
	}

	kind := []byte(arguments[1].String())

	switch len(kind) {
	case 0:
		return errors.New("empty type")
	case 1:
	default:
		return errors.New("excessive type")
	}

	if count < 3 {
		return errors.New("help/description required")
	}

	help := arguments[2].String()

	parameters := &registry.Parameters{
		Name:   name,
		Type:   kind[0],
		Help:   help,
	}

	if count >= 4 {
		labels := make([]string, 0)

		for _, value := range arguments[3].Array() {
			labels = append(labels, value.String())
		}

		if len(labels) > 0 {
			parameters.Labels = labels
		}
	}

	switch parameters.Type {
	case types.Histogram:
		parameters.Histogram = struct {
			Buckets []float64
		}{}
	case types.Summary:
		parameters.Summary = struct {
			Objectives map[float64]float64
			Age        time.Duration
			Buckets    uint32
			BufCap     uint32
		}{}
	}

	if count > 4 {
		switch parameters.Type {
		case types.Histogram:
			buckets := make([]float64, 0)

			for _, value := range arguments[4].Array() {
				buckets = append(buckets, value.Float())
			}

			if len(buckets) > 0 {
				parameters.Histogram.Buckets = buckets
			}
		case types.Summary:
			objectives := make(map[float64]float64)

			for i, value := range arguments[4].Array() {
				tuple := value.Array()

				if len(tuple) != 2 {
					return fmt.Errorf("invalid objective %+v at %d", tuple, i)
				}

				objectives[tuple[0].Float()] = tuple[1].Float()
			}

			if len(objectives) > 0 {
				parameters.Summary.Objectives = objectives
			}

			if count < 6 {
				break
			}

			parameters.Summary.Age = time.Duration(arguments[5].Integer())

			if count < 7 {
				break
			}

			parameters.Summary.Buckets = uint32(arguments[6].Integer())

			if count < 8 {
				break
			}

			parameters.Summary.BufCap = uint32(arguments[7].Integer())
		}
	}

	return h.local.Register(parameters)
}

func (h *Handler) record(metric byte, arguments []resp.Value) error {
	count := len(arguments)

	if count < 1 {
		return errors.New("name required")
	}

	name := arguments[0].String()

	if count < 2 {
		return errors.New("value required")
	}

	value := arguments[1].Float()
	labels := make(map[string]string)

	if count == 3 {
		for i, v := range arguments[2].Array() {
			tuple := v.Array()

			if len(tuple) != 2 {
				return fmt.Errorf("invalid label tuple at %d", i)
			}

			labels[tuple[0].String()] = tuple[1].String()
		}
	}

	switch metric {
	case types.Histogram:
		histogram, err := h.local.Histogram(name)

		if err != nil {
			return err
		}

		histogram.With(labels).Observe(value)
	case types.Counter:
		counter, err := h.local.Counter(name)

		if err != nil {
			return err
		}

		counter.With(labels).Add(value)
	case types.Summary:
		summary, err := h.local.Summary(name)

		if err != nil {
			return err
		}

		summary.With(labels).Observe(value)
	case types.Gauge:
		gauge, err := h.local.Gauge(name)

		if err != nil {
			return err
		}

		gauge.With(labels).Add(value)
	}

	return nil
}

func (h *Handler) metrics() (string, error) {
	return h.local.Metrics()
}
