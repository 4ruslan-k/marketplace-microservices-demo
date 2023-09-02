package cart_test

import (
	cartEntity "cart/internal/domain/entities/cart"
	"testing"

	"gotest.tools/v3/assert"
)

func TestCartReadModel_NewCartReadModel(t *testing.T) {
	t.Parallel()
	type args struct {
		CustomerID string
		Products   []cartEntity.CartReadModelProduct
	}

	testCases := []struct {
		name string
		args args
		want cartEntity.CartReadModel
	}{
		{
			name: "empty_cart",
			args: args{
				CustomerID: "1",
				Products:   nil,
			},
			want: cartEntity.CartReadModel{
				Type:       "cart",
				CustomerID: "1",
				Products:   nil,
				TotalPrice: 0,
			},
		},
		{
			name: "one_product",
			args: args{
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

			want: cartEntity.CartReadModel{
				Type:       "cart",
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
			args: args{
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

			want: cartEntity.CartReadModel{
				Type:       "cart",
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
			args: args{
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

			want: cartEntity.CartReadModel{
				Type:       "cart",
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
			args: args{
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

			want: cartEntity.CartReadModel{
				Type:       "cart",
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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cart := cartEntity.NewCartReadModel(tc.args.CustomerID, tc.args.Products)
			assert.DeepEqual(t, tc.want, cart)
		})
	}
}
