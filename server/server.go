package server

import (
	"errors"
	"fmt"
	"github.com/chaseisabelle/goresp"
	"github.com/chaseisabelle/phprom/command"
	"github.com/chaseisabelle/resptcp"
)

type Server resptcp.Server

func New(host string) *Server {
	return (*Server)(resptcp.New(host, func(values []goresp.Value, err error) ([]goresp.Value, error) {
		if err != nil {
			return []goresp.Value{goresp.NewError(err)}, nil
		}

		if len(values) < 0 {
			err := errors.New("no command specified")

			return []goresp.Value{goresp.NewError(err)}, nil
		}

		bs, err := values[0].Bytes()

		if err == nil && len(bs) != 1 {
			err = fmt.Errorf("malformed command: %s", string(bs))
		}

		if err != nil {
			return []goresp.Value{goresp.NewError(err)}, nil
		}

		var c command.Command

		switch bs[0] {
		case 'R':
		case 'M':

		default:
			err = fmt.Errorf("invalid command: %s", string(bs))
		}

		if err != nil {
			return []goresp.Value{goresp.NewError(err)}, nil
		}

	}, '\000'))
}

func (s *Server) Serve() error {
	server := resptcp.Server(*s)

	return server.Start()
}
