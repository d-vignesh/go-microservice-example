package main

import (
	"fmt"
	"net"
	"os"

	"github.com/hashicorp/go-hclog"
	protos "github.com/d-vignesh/go-microservice-example/currency/protos/currency"
	"github.com/d-vignesh/go-microservice-example/currency/server"
	"github.com/d-vignesh/go-microservice-example/currency/data"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()

	rates, err := data.NewRates(log)
	if err != nil {
		log.Error("unable to generate rates", "error", err)
		os.Exit(1)
	}

	// create a new gRPC server, use WithInsecure to allow http connections
	gs := grpc.NewServer()

	// create an instance of the currency server
	c := server.NewCurrency(rates, log)

	// register the currency server
	protos.RegisterCurrencyServer(gs, c)

	// register the reflection service which allow clients to determine the methods
	// for this gRPC service
	reflection.Register(gs)

	// create a TCP socket for inbound server connections
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", 9092))
	if err != nil {
		log.Error("unable to create listener", "error", err)
		os.Exit(1)
	}

	// listen for requests
	gs.Serve(l)
}