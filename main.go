package main

import (
	"flag"
	"github.com/chaseisabelle/phprom/configuration"
	"github.com/chaseisabelle/phprom/environment"
	"github.com/chaseisabelle/phprom/server"
)

func main() {
	file := flag.String("config", "config.yml", "the config file path")

	flag.Parse()

	cfg, err := configuration.New(*file)

	if err != nil {
		panic(err)
	}

	env := environment.New(cfg)
	srv := server.New(env)

	err = srv.Serve()

	if err != nil {
		panic(err)
	}
}