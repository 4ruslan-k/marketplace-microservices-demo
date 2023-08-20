package dto

import (
	"catalog_service/internal/domain/entities/product"
)

type ProductListOutput struct {
	Type string          `json:"type"`
	Data []ProductOutput `json:"data"`
}

func NewProductListOutputFromEntities(products []product.Product) ProductListOutput {
	productOutputList := make([]ProductOutput, 0, len(products))
	for _, product := range products {
		productOutputList = append(productOutputList, NewProductOutputFromEntity(product))
	}
	return ProductListOutput{Type: "list", Data: productOutputList}
}
