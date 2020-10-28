// Package classification Product API
//
// Documentation for Product API
//
// 	Schemes: http
// 	BasePath: /
// 	version: 1.0.0
//
//	Consumes:
//	- application/json
//
//  Produces:
//	- application/json
//
// swagger:meta
package handlers

import "github.com/d-vignesh/go-microservice-example/product-api/data"

//
// NOTE: Types defined here are purely for documentation purposes
// these types are not used by any of the handlers

// Generic error message returned as a string
// swagger:response errorResponse
type errorResponseWrapper struct {
	// Description of the error
	// in: body
	Body GenericError
}

// Validation errors defined as an array of strings
// swagger:response errorValidation
type errorValidationWrapper struct {
	// collection of errors
	// in: body
	Body ValidationError
}

// A list of products
// swagger:response productsResponse
type productsResponseWrapper struct {
	// all current products
	// in: body
	Body []data.Product
}

// Data structure representing a single product
// swagger:response productResponse
type productResponseWrapper struct {
	// newly created product
	// in: body
	Body data.Product
}

// no content is returned by this API endpoint
// swagger:response noContentResponse
type noContentResponseWrapper struct {
}

// swagger:parameters updateProduct createProduct
type productParamsWrapper struct {
	// Product data structure to update or create.
	// Note: the id field is ignored by update and create operations
	// in: body
	// required: true
	Body data.Product
}

// swagger:parameters listSingleProduct deleteProduct
type productIDParamsWrapper struct {
	// the ID of the product for which the operation relates
	// in: path
	// required: true
	ID int `json:"id"`
}

// swagger:parameters listProducts listSingleProduct
type productQueryParam struct {
	// Currency used when returning the price of the product,
	// when not specified currency is returned in GBP.
	// in: query
	// required: false
	Currency string 
}