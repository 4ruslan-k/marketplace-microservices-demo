package repositories

import (
	productEntity "cart/internal/domain/entities/product"
	"context"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product productEntity.Product) error
	GetProductByID(ctx context.Context, id string) (productEntity.Product, error)
	UpdateProductByID(ctx context.Context, updatedProduct productEntity.Product) error
	DeleteProductByID(ctx context.Context, id string) error
}
