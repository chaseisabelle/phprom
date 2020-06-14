package server

import (
	"errors"
	"fmt"
	"github.com/chaseisabelle/phprom/handler"
	"github.com/chaseisabelle/resptcp"
	"github.com/tidwall/resp"
)

type Server resptcp.Server

func New(host string, handler *handler.Handler) *Server {
	return (*Server)(resptcp.New(host, func(value resp.Value, err error) (resp.Value, error) {
		if err != nil {
			return resp.ErrorValue(err), nil
		}

		commands := value.Array()

		if value.Type() != resp.Array {
			err = fmt.Errorf("invalid command array %+v", value)
		}

		if err == nil && len(commands) == 0 {
			err = errors.New("empty command array")
		}

		if err != nil {
			return resp.ErrorValue(err), nil
		}

		for i, command := range commands {
			switch command.Type() {
			case resp.SimpleString:
			case resp.Array:
				break
			default:
				return resp.ErrorValue(fmt.Errorf("invalid command %+v at %d", command, i)), nil
			}
		}

		command := commands[0]

		if command.Type() != resp.SimpleString {
			return resp.ErrorValue(fmt.Errorf("invalid main command %+v", command)), nil
		}

		toBytes := []byte(command.String())

		if len(toBytes) == 0 {

		}

		switch len(toBytes) {
		case 0:
			return resp.ErrorValue(errors.New("empty main command")), nil
		case 1:
		default:
			return resp.ErrorValue(errors.New("excessive main command")), nil
		}

		return handler.Handle(toBytes[0], commands[1:]), nil
	}))
}

func (s *Server) Serve() error {
	server := resptcp.Server(*s)

	return server.Start()
}