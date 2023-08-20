package cart

import (
	customErrors "shared/errors"
)

var (
	ErrInvalidCustomerID = customErrors.NewIncorrectInputError("cart.products.add.invalid_customer_id", "invalid customer ID")
)

// TODO: add created_at/updated_at
type CartProduct struct {
	ProductID string
	Quantity  int
}

type Cart struct {
	customerID string
	products   []CartProduct
	// TODO: test
	events []Event
}

type AddedProduct struct {
	Product CartProduct
}

func (a AddedProduct) EventType() string {
	return "added_product"
}

type ProductQuantityChanged struct {
	Product CartProduct
}

func (p ProductQuantityChanged) EventType() string {
	return "product_quantity_changed"
}

type ProductRemoved struct {
	ProductID string
}

func (r ProductRemoved) EventType() string {
	return "product_removed"
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

func (c Cart) Events() []Event {
	return c.events
}

type Event interface {
	EventType() string
}

func (cart Cart) UpdateProductsInCart(productToUpdate CartProduct) (Cart, error) {
	if productToUpdate.Quantity <= 0 {
		return cart.deleteProductFromCart(productToUpdate)
	}
	isProductInCart := false
	for i := range cart.products {
		productInCart := &cart.products[i]
		if productInCart.ProductID == productToUpdate.ProductID {
			isProductInCart = true
			productInCart.Quantity = productToUpdate.Quantity
			cart.events = append(cart.events, ProductQuantityChanged{Product: *productInCart})
		}
	}
	if !isProductInCart {
		cart.products = append(cart.products, productToUpdate)
		cart.events = append(cart.events, AddedProduct{Product: productToUpdate})
	}
	return cart, nil
}

func (cart Cart) deleteProductFromCart(productToRemove CartProduct) (Cart, error) {
	isProductInCart := false
	var indexOfProductToDelete int
	for i := range cart.products {
		productInCart := &cart.products[i]
		if productInCart.ProductID == productToRemove.ProductID {
			isProductInCart = true
			indexOfProductToDelete = i
		}
	}

	if isProductInCart {
		cart.products = append(cart.products[:indexOfProductToDelete], cart.products[indexOfProductToDelete+1:]...)
		cart.events = append(cart.events, ProductRemoved{ProductID: productToRemove.ProductID})
	}
	return cart, nil
}

func (cart Cart) isZero() bool {
	return cart.customerID == ""
}
