package command

import (
	"errors"
	"fmt"
	"github.com/chaseisabelle/goresp"
	"github.com/chaseisabelle/phprom/environment"
)

const (
	METRICS  = 'M'
	REGISTER = 'R'
	RECORD   = 'C'
	CLOSE    = 'X'
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
	case RECORD:
		cmd = NewRecord(e)
	case CLOSE:
		cmd = NewClose(e)
	default:
		err = fmt.Errorf("invalid command: %s", string(id))
	}

	return cmd, err
}

func mapLabels(vals []goresp.Value) (map[string]string, error) {
	labs := make(map[string]string)

	for _, val := range vals {
		tup, err := val.Array()

		if err != nil {
			return nil, err
		}

		if len(tup) != 2 {
			return nil, errors.New("invalid tuple")
		}

		key, err := tup[0].String()

		if err != nil {
			return nil, err
		}

		lab, err := tup[1].String()

		if err != nil {
			return nil, err
		}

		labs[key] = lab
	}

	return labs, nil
}
