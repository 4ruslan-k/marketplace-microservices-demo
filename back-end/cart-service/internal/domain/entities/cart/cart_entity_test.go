package cart_test

import (
	cartEntity "cart_service/internal/domain/entities/cart"
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
	type expRes struct {
		CustomerID string
		Products   []cartEntity.CartProduct
	}
	testCases := []struct {
		name   string
		in     cartEntity.CreateCartParams
		expRes expRes
		expErr error
	}{
		{
			name: "err_invalid_customer_id",
			in: cartEntity.CreateCartParams{
				CustomerID: "",
			},
			expErr: cartEntity.ErrInvalidCustomerID,
		},
		{
			name: "empty_cart",
			in: cartEntity.CreateCartParams{
				CustomerID: "123",
			},
			expRes: expRes{
				CustomerID: "123",
				Products:   nil,
			},
			expErr: nil,
		},
		{
			name: "cart_with_products",
			in: cartEntity.CreateCartParams{
				CustomerID: "123",
				Products: []cartEntity.CartProduct{
					{
						ProductID: "123",
						Quantity:  10,
					},
				},
			},
			expRes: expRes{
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
		t.Run(tc.name, func(t *testing.T) {
			cart, err := cartEntity.NewCart(tc.in)

			require.Equal(t, tc.expErr, err)
			assert.Equal(t, tc.expRes.CustomerID, cart.CustomerID())
			assert.DeepEqual(t, tc.expRes.Products, cart.Products())
		})
	}
}

func TestCartEntity_UpdateProductsInCart(t *testing.T) {
	type expRes struct {
		CustomerID string
		Products   []cartEntity.CartProduct
		Events     []cartEntity.Event
	}
	testCases := []struct {
		name   string
		in     cartEntity.CartProduct
		cart   cartEntity.Cart
		expRes expRes
		expErr error
	}{
		{
			name: "add_product_to_empty_cart",
			cart: NewTestCart(t, "123", nil),
			in: cartEntity.CartProduct{
				ProductID: "123",
				Quantity:  2,
			},
			expRes: expRes{
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
			in: cartEntity.CartProduct{
				ProductID: "product_id_two",
				Quantity:  2,
			},
			expRes: expRes{
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
			in: cartEntity.CartProduct{
				ProductID: "product_id_one",
				Quantity:  2,
			},
			expRes: expRes{
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
			in: cartEntity.CartProduct{
				ProductID: "product_id_one",
				Quantity:  0,
			},
			expRes: expRes{
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
			in: cartEntity.CartProduct{
				ProductID: "123",
				Quantity:  0,
			},
			expRes: expRes{
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
			in: cartEntity.CartProduct{
				ProductID: "product_id_two",
				Quantity:  0,
			},
			expRes: expRes{
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
			in: cartEntity.CartProduct{
				ProductID: "product_id_one",
				Quantity:  0,
			},
			expRes: expRes{
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
		t.Run(tc.name, func(t *testing.T) {
			cart, err := tc.cart.UpdateProductsInCart(tc.in)

			require.Equal(t, tc.expErr, err)
			assert.Equal(t, tc.expRes.CustomerID, cart.CustomerID())
			assert.DeepEqual(t, tc.expRes.Products, cart.Products())
			assert.DeepEqual(t, tc.expRes.Events, cart.Events())
		})
	}
}
