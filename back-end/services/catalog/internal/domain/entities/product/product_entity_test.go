package product_test

import (
	"catalog/internal/domain/entities/product"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProduct(t *testing.T) {
	t.Parallel()
	type want struct {
		Name      string
		Price     float64
		Quantity  int
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	testCases := []struct {
		name   string
		args   product.CreateProductParams
		want   want
		expErr error
	}{
		{
			name: "ValidParams_ReturnsProduct",
			args: product.CreateProductParams{
				Name:     "Test Product",
				Price:    9.99,
				Quantity: 10,
			},
			want: want{
				Name:      "Test Product",
				Price:     9.99,
				Quantity:  10,
				CreatedAt: time.Now(),
			},
			expErr: nil,
		},
		{
			name: "InvalidName_ReturnsError",
			args: product.CreateProductParams{
				Name:     "",
				Price:    0,
				Quantity: 10,
			},
			want:   want{},
			expErr: product.ErrInvalidProductName,
		},
		{
			name: "InvalidPrice_ReturnsError",
			args: product.CreateProductParams{
				Name:     "Test Product",
				Price:    0,
				Quantity: 10,
			},
			want:   want{},
			expErr: product.ErrInvalidProductPrice,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			p, err := product.NewProduct(tc.args)

			require.Equal(t, tc.expErr, err)
			assert.Equal(t, tc.want.Name, p.Name())
			assert.Equal(t, tc.want.Price, p.Price())
			assert.Equal(t, tc.want.Quantity, p.Quantity())
			maxDelta := 2 * time.Millisecond
			assert.True(t, tc.want.CreatedAt.Sub(p.CreatedAt()) <= maxDelta)
		})
	}
}

func TestNewProductFromDatabase(t *testing.T) {
	t.Parallel()
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
