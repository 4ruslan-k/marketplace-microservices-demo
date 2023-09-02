package cart_test

import (
	cartEntity "cart/internal/domain/entities/cart"
	"testing"

	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

func NewTestCart(t *testing.T, customerID string, products []cartEntity.CartProduct) cartEntity.Cart {
	cart, err := cartEntity.NewCart(cartEntity.CreateCartParams{
		CustomerID: customerID,
		Products:   products,
	})
	if err != nil {
		t.Fatal(err)
	}
	return cart
}

func TestCartEntity_NewCart(t *testing.T) {
	t.Parallel()
	type want struct {
		CustomerID string
		Products   []cartEntity.CartProduct
	}
	testCases := []struct {
		name   string
		args   cartEntity.CreateCartParams
		want   want
		expErr error
	}{
		{
			name: "error_invalid_customer_id",
			args: cartEntity.CreateCartParams{
				CustomerID: "",
			},
			expErr: cartEntity.ErrInvalidCustomerID,
		},
		{
			name: "empty_cart",
			args: cartEntity.CreateCartParams{
				CustomerID: "123",
			},
			want: want{
				CustomerID: "123",
				Products:   nil,
			},
			expErr: nil,
		},
		{
			name: "cart_with_products",
			args: cartEntity.CreateCartParams{
				CustomerID: "123",
				Products: []cartEntity.CartProduct{
					{
						ProductID: "123",
						Quantity:  10,
					},
				},
			},
			want: want{
				CustomerID: "123",
				Products: []cartEntity.CartProduct{
					{
						ProductID: "123",
						Quantity:  10,
					},
				},
			},
			expErr: nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cart, err := cartEntity.NewCart(tc.args)

			require.Equal(t, tc.expErr, err)
			assert.Equal(t, tc.want.CustomerID, cart.CustomerID())
			assert.DeepEqual(t, tc.want.Products, cart.Products())
		})
	}
}

func TestCartEntity_UpdateProductsInCart(t *testing.T) {
	t.Parallel()
	type want struct {
		CustomerID string
		Products   []cartEntity.CartProduct
		Events     []cartEntity.Event
	}
	testCases := []struct {
		name   string
		args   cartEntity.CartProduct
		cart   cartEntity.Cart
		want   want
		expErr error
	}{
		{
			name: "add_product_to_empty_cart",
			cart: NewTestCart(t, "123", nil),
			args: cartEntity.CartProduct{
				ProductID: "123",
				Quantity:  2,
			},
			want: want{
				CustomerID: "123",
				Products: []cartEntity.CartProduct{
					{
						ProductID: "123",
						Quantity:  2,
					},
				},
				Events: []cartEntity.Event{cartEntity.AddedProduct{
					Product: cartEntity.CartProduct{
						ProductID: "123",
						Quantity:  2,
					},
				},
				},
			},
			expErr: nil,
		},
		{
			name: "add_product_to_cart_with_another_product",
			cart: NewTestCart(t, "customer_id", []cartEntity.CartProduct{
				{
					ProductID: "product_id_one",
					Quantity:  5,
				},
			}),
			args: cartEntity.CartProduct{
				ProductID: "product_id_two",
				Quantity:  2,
			},
			want: want{
				CustomerID: "customer_id",
				Products: []cartEntity.CartProduct{
					{
						ProductID: "product_id_one",
						Quantity:  5,
					},
					{
						ProductID: "product_id_two",
						Quantity:  2,
					},
				},
				Events: []cartEntity.Event{cartEntity.AddedProduct{
					Product: cartEntity.CartProduct{
						ProductID: "product_id_two",
						Quantity:  2,
					},
				}},
			},
			expErr: nil,
		},
		{
			name: "update_product_quantity",
			cart: NewTestCart(t, "customer_id", []cartEntity.CartProduct{
				{
					ProductID: "product_id_one",
					Quantity:  5,
				},
			}),
			args: cartEntity.CartProduct{
				ProductID: "product_id_one",
				Quantity:  2,
			},
			want: want{
				CustomerID: "customer_id",
				Products: []cartEntity.CartProduct{
					{
						ProductID: "product_id_one",
						Quantity:  2,
					},
				},
				Events: []cartEntity.Event{cartEntity.ProductQuantityChanged{
					Product: cartEntity.CartProduct{
						ProductID: "product_id_one",
						Quantity:  2,
					}},
				},
			},
			expErr: nil,
		},
		{
			name: "delete_product_from_cart",
			cart: NewTestCart(t, "customer_id", []cartEntity.CartProduct{
				{
					ProductID: "product_id_one",
					Quantity:  5,
				},
			}),
			args: cartEntity.CartProduct{
				ProductID: "product_id_one",
				Quantity:  0,
			},
			want: want{
				CustomerID: "customer_id",
				Products:   []cartEntity.CartProduct{},
				Events: []cartEntity.Event{
					cartEntity.ProductRemoved{
						ProductID: "product_id_one",
					},
				},
			},
			expErr: nil,
		},
		{
			name: "delete_product_from_empty_cart",
			cart: NewTestCart(t, "123", nil),
			args: cartEntity.CartProduct{
				ProductID: "123",
				Quantity:  0,
			},
			want: want{
				CustomerID: "123",
				Products:   nil,
			},
			expErr: nil,
		},
		{
			name: "delete_not_added_product_from_cart",
			cart: NewTestCart(t, "customer_id", []cartEntity.CartProduct{
				{
					ProductID: "product_id_one",
					Quantity:  5,
				},
			}),
			args: cartEntity.CartProduct{
				ProductID: "product_id_two",
				Quantity:  0,
			},
			want: want{
				CustomerID: "customer_id",
				Products: []cartEntity.CartProduct{
					{
						ProductID: "product_id_one",
						Quantity:  5,
					},
				},
			},
			expErr: nil,
		},
		{
			name: "delete_product_from_cart_with_the_same_product",
			cart: NewTestCart(t, "customer_id", []cartEntity.CartProduct{
				{
					ProductID: "product_id_one",
					Quantity:  5,
				},
			}),
			args: cartEntity.CartProduct{
				ProductID: "product_id_one",
				Quantity:  0,
			},
			want: want{
				CustomerID: "customer_id",
				Products:   []cartEntity.CartProduct{},
				Events: []cartEntity.Event{
					cartEntity.ProductRemoved{
						ProductID: "product_id_one",
					},
				},
			},
			expErr: nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cart, err := tc.cart.UpdateProductsInCart(tc.args)

			require.Equal(t, tc.expErr, err)
			assert.Equal(t, tc.want.CustomerID, cart.CustomerID())
			assert.DeepEqual(t, tc.want.Products, cart.Products())
			assert.DeepEqual(t, tc.want.Events, cart.Events())
		})
	}
}
