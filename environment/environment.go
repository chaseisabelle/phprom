package environment

import (
	"github.com/chaseisabelle/phprom/configuration"
	"sync"
)

type Environment struct {
	sync.Mutex
	config *configuration.Configuration
}

func New(cfg *configuration.Configuration) *Environment {
	return &Environment{
		config: cfg,
	}
}

func (e *Environment) Host() string {
	e.Lock()

	defer e.Unlock()

	return e.config.Host
}
