package metrics

import (
	"fmt"
)

type Metric interface {
	Type() byte
	Register() error
	Record(float64, map[string]string) error
}

func canRegisterAs(n string, t byte) error {
	m, _ := Registered(n)

	if m == nil {
		return nil
	}

	if m.Type() != t {
		return fmt.Errorf("metric %s registered as %s, not %s", n, typeName(m.Type()), typeName(t))
	}

	return nil
}

func typeName(t byte) string {
	var n string

	switch t {
	case HISTOGRAM:
		n = "histogram"
	case COUNTER:
		n = "counter"
	case GAUGE:
		n = "gauge"
	case SUMMARY:
		n = "summary"
	default:
		n = "?"
	}

	return n
}