package cart

import (
	customErrors "cart_service/pkg/errors"
)

var (
	ErrInvalidCustomerID = customErrors.NewIncorrectInputError("cart.products.add.invalid_customer_id", "invalid customer ID")
)

type CartProduct struct {
	ProductID string
	Quantity  int
}

type Cart struct {
	customerID string
	products   []CartProduct
}

type CreateCartParams struct {
	CustomerID string
	Products   []CartProduct
}

func NewCart(createCartParams CreateCartParams) (Cart, error) {
	if createCartParams.CustomerID == "" {
		return Cart{}, ErrInvalidCustomerID
	}
	cart := Cart{
		customerID: createCartParams.CustomerID,
		products:   createCartParams.Products,
	}
	return cart, nil
}

func (c Cart) CustomerID() string {
	return c.customerID
}

func (c Cart) Products() []CartProduct {
	return c.products
}

func (cart Cart) AddProductToCart(productToAdd CartProduct) (Cart, error) {
	isProductInCart := false
	for _, productInCart := range cart.products {
		if productInCart.ProductID == productToAdd.ProductID {
			isProductInCart = true
			productInCart.Quantity += productToAdd.Quantity
		}
	}
	if !isProductInCart {
		cart.products = append(cart.products, productToAdd)
	}
	return cart, nil
}
