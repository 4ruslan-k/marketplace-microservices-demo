package dto

type CreateProductInput struct {
	Name string `json:"name" binding:"required"`
}
