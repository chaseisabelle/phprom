package main

import (
	"flag"
	"github.com/chaseisabelle/phprom/srv/v1"
	"log"
)

func main() {
	adr := flag.String("address", "0.0.0.0:3333", "the host:port to listen on")
	api := flag.String("api", string(v1.GrpcApi), "the api to use (grpc or rest)")

	flag.Parse()

	srv, err := v1.New(v1.API(*api), *adr)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("listening on " + *adr)

	err = srv.Serve()

	if err != nil {
		log.Fatal(err)
	}
}
