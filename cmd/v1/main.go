package main

import (
	"flag"
	phprom_v1 "github.com/chaseisabelle/phprom/pkg/v1"
	"github.com/chaseisabelle/phprom/src/v1"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

func main() {
	adr := flag.String("address", "0.0.0.0:3333", "the host:port to listen on")

	flag.Parse()

	ins, err := v1.New()

	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", *adr)

	if err != nil {
		log.Fatal(err)
	}

	srv := grpc.NewServer()

	phprom_v1.RegisterServiceServer(srv, ins)

	wg := sync.WaitGroup{}

	log.Println("starting server...")

	wg.Add(1)

	go func() {
		err = srv.Serve(lis)

		if err != nil {
			log.Fatal(err)
		}

		wg.Done()
	}()

	log.Println("server started")
	log.Println("listening on " + *adr)

	wg.Wait()

	log.Println("stopping server...")

	srv.GracefulStop()

	log.Println("server stopped")
	log.Println("i'll be back")
}
