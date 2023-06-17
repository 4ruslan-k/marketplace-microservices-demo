package repositories

import (
	cartEntity "cart_service/internal/domain/entities/cart"
	"context"
)

type CartRepository interface {
	GetByCustomerID(ctx context.Context, customerID string) (cartEntity.CartReadModel, error)
}
