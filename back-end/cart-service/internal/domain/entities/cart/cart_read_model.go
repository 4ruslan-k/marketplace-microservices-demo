package cart

import (
	"math"
)

type CartReadModelProduct struct {
	ProductID string
	Name      string
	Quantity  int
	Price     float64
}

type CartReadModel struct {
	CustomerID string
	Products   []CartReadModelProduct
	TotalPrice float64
}

func NewCartReadModel(customerID string, products []CartReadModelProduct) CartReadModel {
	cart := CartReadModel{
		CustomerID: customerID,
		Products:   products,
	}
	cart.TotalPrice = cart.calculateTotalPrice()
	return cart
}

func (c CartReadModel) calculateTotalPrice() float64 {
	var totalPrice float64
	for _, product := range c.Products {
		totalPrice += product.Price * float64(product.Quantity)
	}

	// format to two decimal places
	totalPriceTwoDecimals := math.Round(totalPrice*100) / 100

	return totalPriceTwoDecimals
}
