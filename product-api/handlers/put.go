package handlers

import (
	"net/http"

	"github.com/d-vignesh/go-microservice-example/product-api/data"
)

// swagger:route PUT /products products updateProduct
// Update a product details
//
// responses:
//		201: noContentResponse
//		404: errorResponse
//		422: errorValidation

// Update handles PUT request to update products
func (p *Products) Update(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	// fetch the product from the context
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	prod.ID = getProductID(r)
	p.l.Debug("updating record id", prod.ID)

	err := p.productDB.UpdateProduct(prod)
	if err == data.ErrProductNotFound {
		p.l.Error("product not found", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: "Product not found in database"}, rw)
		return
	}

	// write the noContent success header
	rw.WriteHeader(http.StatusNoContent)
}