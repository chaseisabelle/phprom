package main

import (
	"flag"
	"github.com/chaseisabelle/phprom/src/v1"
	"github.com/chaseisabelle/stop"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	adr := flag.String("address", "127.0.0.1:3333", "the host:port to listen on")

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

	v1.RegisterServiceServer(srv, ins)

	log.Println("starting server...")

	go func() {
		err = srv.Serve(lis)

		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("server started")

	for !stop.Stopped() {
	}

	log.Println("stopping server...")
	srv.GracefulStop()
	log.Println("server stopped")
	log.Println("i'll be back")
}
