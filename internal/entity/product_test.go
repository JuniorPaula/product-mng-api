package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProduct(t *testing.T) {
	product, err := NewProduct("Product 1", 10.00)
	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.NotEmpty(t, product.ID)
	assert.NotEmpty(t, product.CreatedAt)
	assert.Equal(t, "Product 1", product.Name)
	assert.Equal(t, 10.00, product.Price)
}

func TestProduct_When_Name_Is_Required(t *testing.T) {
	product, err := NewProduct("", 10.00)
	assert.Nil(t, product)
	assert.Equal(t, ErrNameRequired, err)
}

func TestProduct_When_Price_Is_Required(t *testing.T) {
	product, err := NewProduct("Product 1", 0.0)
	assert.Nil(t, product)
	assert.Equal(t, ErrPriceRequired, err)
}

func TestProduct_When_Price_Is_Invalid(t *testing.T) {
	product, err := NewProduct("Product 1", -10.0)
	assert.Nil(t, product)
	assert.Equal(t, ErrPriceInvalid, err)
}

func TestProduct_Validate(t *testing.T) {
	product, err := NewProduct("Product 1", 10.00)
	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.Nil(t, product.Validate())
}
