package main

import (
	"net/http"
	"log"
	"os"
	"os/signal"
	"time"
	"context"

	"github.com/d-vignesh/go-microservice-example/product-api/handlers"
	"github.com/d-vignesh/go-microservice-example/product-api/data"
	protos "github.com/d-vignesh/go-microservice-example/currency/protos/currency"

	"github.com/gorilla/mux"
	gohandlers "github.com/gorilla/handlers"
	"github.com/go-openapi/runtime/middleware"
	"google.golang.org/grpc"
	"github.com/hashicorp/go-hclog"
)


func main() {
	l := hclog.Default()
	v := data.NewValidation()

	conn, err := grpc.Dial("localhost:9092", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// create currency client
	cc := protos.NewCurrencyClient(conn)

	// create a database instance
	db := data.NewProductsDB(cc, l)

	// create the product handler
	ph := handlers.NewProducts(l, v, db)

	// create a new serve mux and register the handler
	sm := mux.NewRouter()

	// handlers for API
	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/products", ph.ListAll).Queries("currency", "{[A-Z]{3}}")
	getR.HandleFunc("/products", ph.ListAll)

	getR.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle).Queries("currency", "{[A-Z]{3}}")
	getR.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle)

	putR := sm.Methods(http.MethodPut).Subrouter()
	putR.HandleFunc("/products/{id:[0-9]+}", ph.Update)
	putR.Use(ph.MiddlewareValidateProduct)

	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/products", ph.Create)
	postR.Use(ph.MiddlewareValidateProduct)

	deleteR := sm.Methods(http.MethodDelete).Subrouter()
	deleteR.HandleFunc("/products/{id:[0-9]+}", ph.Delete)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getR.Handle("/docs", sh)
	getR.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))


	// CORS
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	// create a server
	s := http.Server{
		Addr:		  ":9090",
		Handler:	  ch(sm),
		ErrorLog: 	  l.StandardLogger(&hclog.StandardLoggerOptions{}),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// start the server
	go func() {
		l.Info("Starting server on port 9090")

		err := s.ListenAndServe()
		if err != nil {
			l.Error("error starting server", "error", err)
			os.Exit(1)
		}
	}()

	// look for interupts for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// block until the signal is received
	sig := <-c 
	log.Println("got signal :", sig)

	//gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30 * time.Second)
	s.Shutdown(ctx)
}	