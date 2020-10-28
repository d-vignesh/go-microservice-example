package handlers

import (
	"net/http"

	"github.com/d-vignesh/go-microservice-example/product-api/data"
)

// swagger:route GET /products products listProducts
// Returns a list of products from the database.
// responses:
//		200: productResponse

// ListAll handles GET requests and returns all current products
func (p *Products) ListAll(rw http.ResponseWriter, r *http.Request) {
	p.l.Debug("get all records")
	rw.Header().Add("Content-Type", "application/json")

	cur := r.URL.Query().Get("currency")

	prods, err := p.productDB.GetProducts(cur)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = data.ToJSON(prods, rw)
	if err != nil {
		p.l.Error("unable to serializing product", "error", err)
	}
}

// swagger:route GET /products/{id} products listSingleProduct
// Returns a single requested product or error is product not found
// responses:
//		200: productResponse
//		404: errorResponse

func (p *Products) ListSingle(rw http.ResponseWriter, r *http.Request) {

	rw.Header().Add("Content-Type", "application/json")

	id := getProductID(r)
	cur := r.URL.Query().Get("currency")

	p.l.Debug("got record", "id", id)

	prod, err := p.productDB.GetProductByID(id, cur)

	switch err {
	case nil :

	case data.ErrProductNotFound:
		p.l.Error("unable to fetch product", "error", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return 
	default:
		p.l.Error("unable to fetch product", "error", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = data.ToJSON(prod, rw)
	if err != nil {
		p.l.Error("unable to serialize product", "error", err)
	}
}