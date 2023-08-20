package domainservices

import (
	customerEntity "cart_service/internal/domain/entities/customer"
	"cart_service/internal/ports/repositories"
	"context"
	"fmt"

	"github.com/rs/zerolog"
)

type CustomerDomainService interface {
	CreateCustomer(ctx context.Context, params customerEntity.CreateCustomerParams) (*customerEntity.Customer, error)
}

type customerService struct {
	logger             zerolog.Logger
	CustomerRepository repositories.CustomerRepository
}

func NewCustomerService(
	logger zerolog.Logger,
	CustomerRepository repositories.CustomerRepository,
) *customerService {
	return &customerService{logger, CustomerRepository}
}

// Creates a user from user input or from social account
func (u customerService) CreateCustomer(
	ctx context.Context,
	createUserParams customerEntity.CreateCustomerParams,
) (*customerEntity.Customer, error) {
	u.logger.Info().Interface("createUserParams", createUserParams).Msg("customerService -> CreateCustomer -> createUserParams")
	newUser, err := customerEntity.NewCustomer(customerEntity.CreateCustomerParams{
		Name:  createUserParams.Name,
		Email: createUserParams.Email,
		ID:    createUserParams.ID,
	})

	if err != nil {
		return nil, fmt.Errorf("customerService -> CreateCustomer - NewCustomer: %w", err)
	}

	err = u.CustomerRepository.Create(ctx, *newUser)

	if err != nil {
		return nil, fmt.Errorf("customerService -> CreateCustomer - CustomerRepository.Create: %w", err)
	}

	return nil, err
}
