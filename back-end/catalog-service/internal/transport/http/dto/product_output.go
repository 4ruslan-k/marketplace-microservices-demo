package dto

import "catalog_service/internal/domain/entities/product"

type ProductOutput struct {
	Type  string  `json:"type"`
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func NewProductOutputFromEntity(product product.Product) ProductOutput {
	return ProductOutput{
		Type:  "product",
		ID:    product.ID(),
		Name:  product.Name(),
		Price: product.Price(),
	}
}
