package server

import (
	"context"
	"io"
	"time"

	"github.com/hashicorp/go-hclog"
	protos "github.com/d-vignesh/go-microservice-example/currency/protos/currency"
	"github.com/d-vignesh/go-microservice-example/currency/data"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Currency is a gRPC server it implements the methods defined by the CurrencyServer interface
type Currency struct {
	rates *data.ExchangeRates
	log hclog.Logger
	subscriptions map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest
}

// NewCurrency create a new Currency server
func NewCurrency(er *data.ExchangeRates, l hclog.Logger) *Currency {
	c := &Currency{er, l, make(map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest)}
	go c.handleUpdates()
	return c 
}

func (c *Currency) handleUpdates() {
	ru := c.rates.MonitorRates(20 * time.Second)
	for range ru {
		c.log.Info("got updated rates")

		// loop over subscribed clients
		for k, v := range c.subscriptions {

			// loop over subscribed rates
			for _, rr := range v {
				r, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
				if err != nil {
					c.log.Error("unable to get updated rate", "base", rr.GetBase().String(), "destination", rr.GetDestination().String())
				}

				// create the response and sent to the client
				err = k.Send(&protos.StreamingRateResponse{
					Message: &protos.StreamingRateResponse_RateResponse{
						RateResponse: &protos.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: r},
					},
				})

				if err != nil {
					c.log.Error("unable to send updated rate", "base", rr.GetBase().String(), "destination", rr.GetDestination().String())
				}
			}
		}
	}
}


// GetRate implements the CurrencyServer GetRate method and returns the currency exchange rate
// for the given currencies.
func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("handle request for GetRate", "base", rr.GetBase(), "dest", rr.GetDestination())

	// Validate that base and destination currency are different
	if rr.Base == rr.Destination {
		// create the grpc error and return to the client
		err := status.Errorf(
			codes.InvalidArgument,
			"Base currency %s can not be same as destination currency %s",
			rr.Base.String(),
			rr.Destination.String(),
		)
		return nil, err
	}

	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}

	return &protos.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: rate}, nil
}

// SubscribeRates implements the gRPC bidirectional streaming method for the server
func (c *Currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error {

	// handle client messages
	for {
		rr, err := src.Recv() // Recv is a blocking method which returns on client data
		// io.EOF signals that the client has closed the connection
		if err == io.EOF {
			c.log.Info("client has closed connection")
			break
		}

		// any other error means the transport between the server and client is unavailable
		if err != nil {
			c.log.Error("unable to read from client", "error", err)
			return err
		}

		c.log.Info("handle client request", "request_base", rr.GetBase(), "request_dest", rr.GetDestination())

		rrs, ok := c.subscriptions[src]
		if !ok {
			rrs = []*protos.RateRequest{}
		}

		// check if already in the subscribe list and return a custom gRPC error
		for _, r := range rrs {
			// if we already have subscribed to this currency return an error
			if r.Base == rr.Base && r.Destination == rr.Destination {
				c.log.Error("subscription already active", "base", rr.Base.String(), "dest", rr.Destination.String())

				grpcError := status.New(codes.InvalidArgument, "subscription already active for rate")
				grpcError, err = grpcError.WithDetails(rr)
				if err != nil {
					c.log.Error("unable to add metadate to error message", "error", err)
					continue
				}

				// can't return error as that will terminate the connection, instead must send an error which
				// can be handled by the client recv stream
				srr := &protos.StreamingRateResponse_Error{Error: grpcError.Proto()}
				src.Send(&protos.StreamingRateResponse{Message: srr})
			}
		}

		rrs = append(rrs , rr)
		c.subscriptions[src] = rrs
	}

	return nil
}