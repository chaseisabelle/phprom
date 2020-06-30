package command

import (
	"errors"
	"fmt"
	"github.com/chaseisabelle/goresp"
	"github.com/chaseisabelle/phprom/environment"
	"github.com/chaseisabelle/phprom/registry"
	"github.com/chaseisabelle/phprom/types"
	"github.com/tidwall/resp"
	"time"
)

type Command interface {
	Execute(...goresp.Value) ([]goresp.Value, error)
}

func NewCommand(e *environment.Environment, id byte) (Command, error) {
	var cmd Command
	var err error

	switch id {
	case REGISTER:
		cmd = NewRegister(e)
	case METRICS:
		cmd = NewMetrics(e)
	default:
		err = fmt.Errorf("invalid command: %s", string(bs))
	}

	return cmd, err
}

const (
	METRICS   byte = 'M'
	REGISTER  byte = 'R'
	HISTOGRAM      = types.Histogram
	COUNTER        = types.Counter
	SUMMARY        = types.Summary
	GAUGE          = types.Gauge
)
