package environment

import (
	"github.com/chaseisabelle/phprom/configuration"
	"github.com/chaseisabelle/phprom/registry"
	"github.com/chaseisabelle/phprom/server"
)

type Environment struct {
	config   *configuration.Configuration
	registry *registry.Registry
	server   *server.Server
}

func New(cfg *configuration.Configuration) *Environment {
	return &Environment{
		config: cfg,
	}
}

func (e *Environment) Registry() (*registry.Registry, error) {
	if e.registry == nil {
		reg, err := registry.New()

		if err != nil {
			return nil, err
		}

		e.registry = reg
	}

	return e.registry, nil
}

func (e *Environment) Server() *server.Server {
	if e.server == nil {
		e.server = server.New(e.config.Host)
	}

	return e.server
}
