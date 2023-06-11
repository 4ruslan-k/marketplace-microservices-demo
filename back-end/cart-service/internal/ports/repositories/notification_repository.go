package repositories

import (
	productEntity "cart_service/internal/domain/entities/product"
	"context"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product productEntity.Product) error
}
