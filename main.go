package main

import (
	"flag"
	"github.com/chaseisabelle/phprom/handler"
	"github.com/chaseisabelle/phprom/registry"
	"github.com/chaseisabelle/phprom/server"
)

func main() {
	host := flag.String("host", ":3333", "host and port to listen on")
	namespace := flag.String("namespace", "", "the metric namespace")

	flag.Parse()

	reg, err := registry.New(*namespace)

	if err != nil {
		panic(err)
	}

	han := handler.New(reg)
	srv := server.New(*host, han)

	go func() {
		for err := range srv.Errors {
			printerr(err)
		}
	}()

	err = srv.Serve()

	if err != nil {
		panic(err)
	}
}

func printerr(err error) {
	println(err.Error())
}

func keys(labs map[string]string) []string {
	keys := make([]string, 0)

	for key := range labs {
		keys = append(keys, key)
	}

	return keys
}