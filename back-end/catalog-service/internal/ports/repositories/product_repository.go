package repositories

import (
	productEntity "catalog_service/internal/domain/entities/product"
	"context"
)

type ProductRepository interface {
	Save(ctx context.Context, product productEntity.Product) error
}
