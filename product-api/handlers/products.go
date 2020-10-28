package handlers

import (
	"net/http"
	"strconv"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"

	"github.com/d-vignesh/go-microservice-example/product-api/data"
)


// KeyProduct is a key used for the Product in the context
type KeyProduct struct{}

// Products handler for getting and updating products
type Products struct {
	l hclog.Logger 
	v *data.Validation
	productDB *data.ProductsDB
}

// NewProducts returns a new product handler with given logger, validator and currencyClient
func NewProducts(l hclog.Logger, v *data.Validation, pdb *data.ProductsDB) *Products {
	return &Products{l, v, pdb}
}

// ErrInvalidProductPath is an error message when the product path is not valid
var ErrInvalidProductPath = fmt.Errorf("invalid path, path should be /products/[id]")

// GenericError is a generic error message returned by server
type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a collection 	of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}

// getProductID returns the product ID from the URL
// Panics if it cannot convert the ID to an integer
// this should never happen as the router ensures that
// this is a valid number
func getProductID(r *http.Request) int {
	// parse the product ID from URL
	vars := mux.Vars(r)

	// convert the id into a integer and return
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}

	return id
}

