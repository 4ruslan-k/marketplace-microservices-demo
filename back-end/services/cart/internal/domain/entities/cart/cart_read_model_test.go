package cart_test

import (
	cartEntity "cart/internal/domain/entities/cart"
	"testing"

	"gotest.tools/v3/assert"
)

func TestCartReadModel_NewCartReadModel(t *testing.T) {

	type in struct {
		CustomerID string
		Products   []cartEntity.CartReadModelProduct
	}

	testCases := []struct {
		name   string
		in     in
		expRes cartEntity.CartReadModel
	}{
		{
			name: "empty_cart",
			in: in{
				CustomerID: "1",
				Products:   nil,
			},
			expRes: cartEntity.CartReadModel{
				CustomerID: "1",
				Products:   nil,
				TotalPrice: 0,
			},
		},
		{
			name: "one_product",
			in: in{
				CustomerID: "1",
				Products: []cartEntity.CartReadModelProduct{
					{
						ProductID: "1",
						Quantity:  10,
						Name:      "test",
						Price:     10,
					},
				},
			},

			expRes: cartEntity.CartReadModel{
				CustomerID: "1",
				Products: []cartEntity.CartReadModelProduct{
					{
						ProductID: "1",
						Quantity:  10,
						Name:      "test",
						Price:     10,
					},
				},
				TotalPrice: 100,
			},
		},
		{
			name: "two_products",
			in: in{
				CustomerID: "1",
				Products: []cartEntity.CartReadModelProduct{
					{
						ProductID: "1",
						Quantity:  10,
						Name:      "apple",
						Price:     10,
					},
					{
						ProductID: "2",
						Quantity:  1,
						Name:      "banana",
						Price:     10.55,
					},
				},
			},

			expRes: cartEntity.CartReadModel{
				CustomerID: "1",
				Products: []cartEntity.CartReadModelProduct{
					{
						ProductID: "1",
						Quantity:  10,
						Name:      "apple",
						Price:     10,
					},
					{
						ProductID: "2",
						Quantity:  1,
						Name:      "banana",
						Price:     10.55,
					},
				},
				TotalPrice: 110.55,
			},
		},
		{
			name: "two_products#2",
			in: in{
				CustomerID: "1",
				Products: []cartEntity.CartReadModelProduct{
					{
						ProductID: "1",
						Quantity:  1,
						Name:      "apple",
						Price:     1.55,
					},
					{
						ProductID: "2",
						Quantity:  1,
						Name:      "banana",
						Price:     10.55,
					},
				},
			},

			expRes: cartEntity.CartReadModel{
				CustomerID: "1",
				Products: []cartEntity.CartReadModelProduct{
					{
						ProductID: "1",
						Quantity:  1,
						Name:      "apple",
						Price:     1.55,
					},
					{
						ProductID: "2",
						Quantity:  1,
						Name:      "banana",
						Price:     10.55,
					},
				},
				TotalPrice: 12.1,
			},
		},
		{
			name: "one_product",
			in: in{
				CustomerID: "1",
				Products: []cartEntity.CartReadModelProduct{
					{
						ProductID: "1",
						Quantity:  3,
						Name:      "apple",
						Price:     1.27,
					},
					{
						ProductID: "2",
						Quantity:  2,
						Name:      "banana",
						Price:     12.53,
					},
				},
			},

			expRes: cartEntity.CartReadModel{
				CustomerID: "1",
				Products: []cartEntity.CartReadModelProduct{
					{
						ProductID: "1",
						Quantity:  3,
						Name:      "apple",
						Price:     1.27,
					},
					{
						ProductID: "2",
						Quantity:  2,
						Name:      "banana",
						Price:     12.53,
					},
				},
				TotalPrice: 28.87,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cart := cartEntity.NewCartReadModel(tc.in.CustomerID, tc.in.Products)
			assert.DeepEqual(t, tc.expRes, cart)
		})
	}
}
