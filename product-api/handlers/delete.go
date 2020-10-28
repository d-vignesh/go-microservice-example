package handlers

import (
	"net/http"

	"github.com/d-vignesh/go-microservice-example/product-api/data"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Delete the product with given ID
//
// responses:
//		201: noContentResponse
//		404: errorResponse
//		501: errorResponse

// Delete handles DELETE request and removes item from database
func (p *Products) Delete(rw http.ResponseWriter, r *http.Request) {
	id := getProductID(r)

	p.l.Debug("deleting record id", id)

	err := p.productDB.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		p.l.Error("product with given id does not exist")

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	if err != nil {
		p.l.Error("unable to delete record", "error", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}