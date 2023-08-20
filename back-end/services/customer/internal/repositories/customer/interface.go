package repository

import (
	"context"
	customerEntity "customer/internal/domain/entities/customer"
)

type CustomerRepository interface {
	GetByID(ctx context.Context, ID string) (*customerEntity.Customer, error)
	GetByEmail(ctx context.Context, email string) (*customerEntity.Customer, error)
	Create(ctx context.Context, customer customerEntity.Customer) error
	Update(ctx context.Context, customer customerEntity.Customer) error
	Delete(ctx context.Context, ID string) error
}
