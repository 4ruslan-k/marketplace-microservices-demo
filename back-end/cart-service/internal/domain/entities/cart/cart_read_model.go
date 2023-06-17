package cart

import (
	"math"
)

type CartReadModelProduct struct {
	ProductID string  `json:"productId"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type CartReadModel struct {
	CustomerID string                 `json:"customerId"`
	Products   []CartReadModelProduct `json:"products"`
	TotalPrice float64                `json:"totalPrice"`
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

func (c CartReadModel) IsZero() bool {
	return c.CustomerID == ""
}
