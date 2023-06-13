package product_test

import (
	"catalog_service/internal/domain/entities/product"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProduct(t *testing.T) {
	type expRes struct {
		Name      string
		Price     float64
		Quantity  int
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	testCases := []struct {
		name   string
		in     product.CreateProductParams
		expRes expRes
		expErr error
	}{
		{
			name: "ValidParams_ReturnsProduct",
			in: product.CreateProductParams{
				Name:     "Test Product",
				Price:    9.99,
				Quantity: 10,
			},
			expRes: expRes{
				Name:      "Test Product",
				Price:     9.99,
				Quantity:  10,
				CreatedAt: time.Now(),
			},
			expErr: nil,
		},
		{
			name: "InvalidName_ReturnsError",
			in: product.CreateProductParams{
				Name:     "",
				Price:    0,
				Quantity: 10,
			},
			expRes: expRes{},
			expErr: product.ErrInvalidProductName,
		},
		{
			name: "InvalidPrice_ReturnsError",
			in: product.CreateProductParams{
				Name:     "Test Product",
				Price:    0,
				Quantity: 10,
			},
			expRes: expRes{},
			expErr: product.ErrInvalidProductPrice,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := product.NewProduct(tc.in)

			require.Equal(t, tc.expErr, err)
			assert.Equal(t, tc.expRes.Name, p.Name())
			assert.Equal(t, tc.expRes.Price, p.Price())
			assert.Equal(t, tc.expRes.Quantity, p.Quantity())
			maxDelta := 2 * time.Millisecond
			assert.True(t, tc.expRes.CreatedAt.Sub(p.CreatedAt()) <= maxDelta)
		})
	}
}

func TestNewProductFromDatabase(t *testing.T) {
	id := "123"
	name := "Test Product"
	price := 9.99
	quantity := 10
	createdAt := time.Now()
	updatedAt := time.Now().Add(1 * time.Hour)

	p := product.NewProductFromDatabase(id, name, price, quantity, createdAt, updatedAt)

	assert.Equal(t, id, p.ID())
	assert.Equal(t, name, p.Name())
	assert.Equal(t, price, p.Price())
	assert.Equal(t, quantity, p.Quantity())
	assert.Equal(t, createdAt, p.CreatedAt())
	assert.Equal(t, updatedAt, p.UpdatedAt())
	assert.False(t, p.IsZero())
}
