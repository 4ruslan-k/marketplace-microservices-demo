package cart

import (
	customErrors "cart_service/pkg/errors"
)

var (
	ErrInvalidCustomerID = customErrors.NewIncorrectInputError("cart.products.add.invalid_customer_id", "invalid customer ID")
)

type ProductAdded struct {
	Product CartProduct
}

type ProductQuantityChanged struct {
	Product CartProduct
}

type ProductRemoved struct {
	ProductID string
}

// TODO: add created_at/updated_at
type CartProduct struct {
	ProductID string
	Quantity  int
}

type Cart struct {
	customerID string
	products   []CartProduct
	_events    []any
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

func (c Cart) Events() []any {
	return c._events
}

func (cart Cart) AddProductToCart(productToAdd CartProduct) (Cart, error) {
	isProductInCart := false
	for i := range cart.products {
		productInCart := &cart.products[i]
		if productInCart.ProductID == productToAdd.ProductID {
			isProductInCart = true
			productInCart.Quantity += productToAdd.Quantity
			cart._events = append(cart._events, ProductQuantityChanged{Product: *productInCart})
		}
	}
	if !isProductInCart {
		cart.products = append(cart.products, productToAdd)
		cart._events = append(cart._events, ProductAdded{Product: productToAdd})
	}
	return cart, nil
}

func (cart Cart) DeleteProductFromCart(productToRemove CartProduct) (Cart, error) {
	shouldBeDeleted := false
	var indexOfProductToDelete int
	for i := range cart.products {
		productInCart := &cart.products[i]
		if productInCart.ProductID == productToRemove.ProductID {
			productInCart.Quantity -= productToRemove.Quantity
			if productInCart.Quantity <= 0 {
				shouldBeDeleted = true
				indexOfProductToDelete = i
			} else {
				cart._events = append(cart._events, ProductQuantityChanged{Product: *productInCart})
			}
		}
	}

	if shouldBeDeleted {
		cart.products = append(cart.products[:indexOfProductToDelete], cart.products[indexOfProductToDelete+1:]...)
		cart._events = append(cart._events, ProductRemoved{ProductID: productToRemove.ProductID})
	}
	return cart, nil
}

func (cart Cart) isZero() bool {
	return cart.customerID == ""
}
