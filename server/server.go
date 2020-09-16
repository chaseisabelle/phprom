package server

import (
	"errors"
	"fmt"
	"github.com/chaseisabelle/goresp"
	"github.com/chaseisabelle/phprom/command"
	"github.com/chaseisabelle/phprom/environment"
	"github.com/chaseisabelle/resptcp"
)

type Server resptcp.Server

func New(env *environment.Environment) *Server {
	return (*Server)(resptcp.New(env.Host(), func(values []goresp.Value, err error) ([]goresp.Value, error) {
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

		cmd, err := command.NewCommand(env, bs[0])

		if err != nil {
			return []goresp.Value{goresp.NewError(err)}, nil
		}

		return cmd.Execute(values[1:]...)
	}, '\000'))
}

func (s *Server) Serve() error {
	server := resptcp.Server(*s)

	go func() {
		for err := range server.Errors {
			println(err.Error())
		}
	}()

	return server.Start()
}
