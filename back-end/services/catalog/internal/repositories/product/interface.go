package repository

import (
	productEntity "catalog/internal/domain/entities/product"
	"context"
)

type ProductRepository interface {
	SaveProduct(ctx context.Context, product productEntity.Product) error
	GetProducts(ctx context.Context) ([]productEntity.Product, error)
	GetProductByID(ctx context.Context, productID string) (productEntity.Product, error)
	DeleteProductByID(ctx context.Context, productID string) error
	UpdateProductByID(ctx context.Context, product productEntity.Product) error
}
