package data

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductMissingNameReturnsErr(t *testing.T) {
	p := Product {
		Price: 1.22,
	}

	v := NewValidation()
	err := v.Validate(p)
	assert.Len(t, err, 1)
}

func TestProductMissingPriceRetrunsErr(t *testing.T) {
	p := Product {
		Name: "abc",
		Price: -1,
	}

	v := NewValidation()
	err := v.Validate(p)
	assert.Len(t, err, 1)
}

func TestProductInvalidSKUReturnsErr(t *testing.T) {
	p := Product {
		Name: "abc",
		Price: 1.22,
		SKU: "abc",
	}

	v := NewValidation()
	err := v.Validate(p)
	assert.Len(t, err, 1)
}

func TestValidProductDoesNotReturnsErr(t *testing.T) {
	p := Product {
		Name: "abc",
		Price: 1.22,
		SKU: "abc-abc-abc",
	}

	v := NewValidation()
	err := v.Validate(p)
	assert.Len(t, err, 0)
}

func TestProductsToJSON(t *testing.T) {
	ps := []*Product{
		&Product {
			Name: "abc",
		},
	}

	b := bytes.NewBufferString("")
	err := ToJSON(ps, b)
	assert.NoError(t, err)
}