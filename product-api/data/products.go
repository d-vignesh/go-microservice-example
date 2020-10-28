package data

import (
	"fmt"
	"time"
	"context"
	"math"

	"github.com/hashicorp/go-hclog"
	protos "github.com/d-vignesh/go-microservice-example/currency/protos/currency"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrProductNotFound is an error raised when a product cannot be found in the database
var ErrProductNotFound = fmt.Errorf("product not found")

// Product defines the structure for an API product
// swagger: model
type Product struct {
	// id for the product
	//
	// required: false
	// min: 1
	ID			int			`json:"id"` // Unique identifier for the product
	
	// name for the product
	//
	// required: true
	// max length: 255
	Name		string		`json:"name" validate:"required"`	

	// description for the product
	//
	// required: false
	// max length: 10000
	Description string		`json:"description"`

	// price of the product
	//
	// required: true
	// min: 0.01
	Price		float64		`json:"price" validate:"gt=0"`

	// SKU for the product
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU			string		`json:"sku" validate:"required,sku"`

	CreatedOn	string		`json:"-"`
	UpdatedOn	string		`json:"-"`
	DeletedOn	string		`json:"-"`
}

// Products defines a slice of Product
type Products []*Product

type ProductsDB struct {
	currency protos.CurrencyClient
	log 	 hclog.Logger
	rates	 map[string]float64
	client   protos.Currency_SubscribeRatesClient
}

func NewProductsDB(c protos.CurrencyClient, l hclog.Logger) *ProductsDB {
	pb := &ProductsDB{c, l, make(map[string]float64), nil}

	go pb.handleUpdates()

	return pb
}

func (p *ProductsDB) handleUpdates() {
	sub, err := p.currency.SubscribeRates(context.Background())
	if err != nil {
		p.log.Error("unable to subscribe for rates", "error", err)
	}
	p.client = sub

	for {
		// Recv returns a StreamingRateResponse which can contain one of two message
		// RateResponse or an Error.
		// We need to handle each case separately
		srr, err := sub.Recv()

		// handle connection errors
		// this is normally terminal required a reconnect
		if err != nil {
			p.log.Error("error receiving message", "error", err)
			return
		}

		// handle a returned error message
		if ge := srr.GetError(); ge != nil {
			sre := status.FromProto(ge)

			if sre.Code() == codes.InvalidArgument {
				errDetails := ""
				if d := sre.Details(); len(d) > 0 {
					p.log.Error("Deets", "d", d)
					if rr, ok := d[0].(*protos.RateRequest); ok {
						errDetails = fmt.Sprintf("base: %s destination: %s", rr.GetBase().String(), rr.GetDestination().String())
					}
				}

				p.log.Error("received error from currency service rate subscription", "error", ge.GetMessage(), "details", errDetails)
			}
		}

		// handle the rate response
		if rr := srr.GetRateResponse(); rr != nil {
			p.log.Info("received updated rate from server", "dest", rr.GetDestination().String())
			p.rates[rr.Destination.String()] = rr.Rate
		}
	}
}

// GetProducts returns all products from the database
func (p *ProductsDB) GetProducts(currency string) (Products, error) {
	if currency == "" {
		return productList, nil
	}

	rate, err := p.getRate(currency)
	if err != nil {
		p.log.Error("unable to get rate", "currency", currency, "error", err)
		return nil, err
	}

	pr := Products{}
	for _, prod := range productList {
		np := *prod
		np.Price = math.Round(np.Price * rate * 10) / 10
		pr = append(pr, &np)
	}
	return pr, nil
}

// GetProductByID returns a single product which matches the id from the database.
// if a product is not found this function returns a ProductNotFound error
func (p *ProductsDB) GetProductByID(id int, currency string) (*Product, error) {
	i := findIndexByProductID(id)
	if i == -1 {
		return nil, ErrProductNotFound
	}

	if currency == "" {
		return productList[i], nil
	}

	rate, err := p.getRate(currency)
	if err != nil {
		p.log.Error("unable to get rate", "currency", currency, "error", err)
		return nil, err
	}

	np := *productList[i]
	np.Price = np.Price * rate 

	return &np, nil
}

// AddProduct adds a new product to the database
func (p *ProductsDB) AddProduct(pr Product) {
	// get the next id in sequence
	maxID := productList[len(productList)-1].ID
	pr.ID = maxID + 1
	productList = append(productList, &pr)
} 

// UpdateProduct replaces a product in the database with the given item
// If a product with the given id does not exists in the database
// this function returns a ProductNotFound error
func (p *ProductsDB) UpdateProduct(pr Product) error {
	i := findIndexByProductID(pr.ID)
	if i == -1 {
		return ErrProductNotFound
	}

	// update the product in the DB
	productList[i] = &pr

	return nil
}

// DeleteProduct deletes a product from the database
func (p *ProductsDB) DeleteProduct(id int) error {
	i := findIndexByProductID(id)
	if i == -1 {
		return ErrProductNotFound
	}
	if i == len(productList) - 1 {
		productList = productList[:i]
	} else {
		productList = append(productList[:i], productList[(i+1):]...)
	}

	return nil
}

// findIndex finds the index of a product in the database
// returns -1 when no product can be found
func findIndexByProductID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}

	return -1
}

func (p *ProductsDB) getRate(destination string) (float64, error) {
	// if cached return
	if r, ok := p.rates[destination]; ok {
		return r, nil
	}

	rr := &protos.RateRequest {
		Base:			protos.Currencies(protos.Currencies_value["EUR"]),
		Destination:	protos.Currencies(protos.Currencies_value[destination]),
	}

	resp, err := p.currency.GetRate(context.Background(), rr)
	if err != nil {
		// convert the grpc error message
		grpcError, ok := status.FromError(err)
		if !ok {
			return -1, err
		}

		// if this is an invalid arguments exception santise the message before returning
		if grpcError.Code() == codes.InvalidArgument {
			return -1, fmt.Errorf("unable to retrieve exchange rate from currency service: %s", grpcError.Message())
		}
	}
	
	p.rates[destination] = resp.Rate

	// subscribe for updates
	p.client.Send(rr)

	return resp.Rate, err
}

// productList is a hard coded list of products for this example data source
var productList = []*Product{
	&Product{
		ID: 		 1,
		Name:		 "Latte",
		Description: "Frothy milky coffee",
		Price:		 2.45,
		SKU:		 "abc323",
		CreatedOn:	 time.Now().UTC().String(),
		UpdatedOn:	 time.Now().UTC().String(),
	},
	&Product{
		ID:			 2,
		Name:		 "Espresso",
		Description: "Short and strong coffee withoud milk",
		Price:		 1.99,
		SKU:		 "fjd34",
		CreatedOn:	 time.Now().UTC().String(),
		UpdatedOn:	 time.Now().UTC().String(),
	},
}

