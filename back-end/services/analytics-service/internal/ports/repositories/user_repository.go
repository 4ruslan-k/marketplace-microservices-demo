package repositories

import (
	userEntity "analytics_service/internal/domain/entities/user"
	"context"
)

type UserRepository interface {
	GetByID(ctx context.Context, ID string) (*userEntity.User, error)
	GetByEmail(ctx context.Context, email string) (*userEntity.User, error)
	Create(ctx context.Context, user userEntity.User) error
	Update(ctx context.Context, user userEntity.User) error
	Delete(ctx context.Context, ID string) error
}
