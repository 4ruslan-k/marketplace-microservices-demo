package repositories

import (
	cartEntity "cart/internal/domain/entities/cart"
	"context"
)

type CartRepository interface {
	GetByCustomerID(ctx context.Context, customerID string) (cartEntity.CartReadModel, error)
	SaveCart(ctx context.Context, customerID string, updateFunc func(cart cartEntity.Cart) (cartEntity.Cart, error)) error
}
