package dto

type CreateProductInput struct {
	Name     string  `json:"name" binding:"required"`
	Price    float64 `json:"price" binding:"required,gte=0"`
	Quantity int     `json:"quantity" binding:"gte=0"`
}
