package dto

type UpdateProductInput struct {
	Name     *string  `json:"name" binding:"omitempty,min=1"`
	Price    *float64 `json:"price" binding:"omitempty,gt=0"`
	Quantity *int     `json:"quantity" binding:"omitempty,gte=0"`
}
