package domainservices

import (
	"context"
	customerEntity "customer_service/internal/domain/entities/customer"
	repository "customer_service/internal/repositories/customer"
	"fmt"

	"github.com/rs/zerolog"
)

type CustomerDomainService interface {
	CreateCustomer(ctx context.Context, params customerEntity.CreateCustomerParams) (*customerEntity.Customer, error)
}

type customerService struct {
	logger             zerolog.Logger
	customerRepository repository.CustomerRepository
}

func NewCustomerService(
	logger zerolog.Logger,
	customerRepository repository.CustomerRepository,
) *customerService {
	return &customerService{logger, customerRepository}
}

// Creates a customer from customer input or from social account
func (u customerService) CreateCustomer(ctx context.Context, createCustomerParams customerEntity.CreateCustomerParams) (*customerEntity.Customer, error) {
	u.logger.Info().Interface("createCustomerParams", createCustomerParams).Msg("customerService -> CreateCustomer -> createCustomerParams")
	newUser, err := customerEntity.NewCustomer(customerEntity.CreateCustomerParams{
		Name:  createCustomerParams.Name,
		Email: createCustomerParams.Email,
		ID:    createCustomerParams.ID,
	})

	if err != nil {
		return nil, fmt.Errorf("customerService -> CreateCustomer - NewCustomer: %w", err)
	}

	err = u.customerRepository.Create(ctx, *newUser)

	if err != nil {
		return nil, fmt.Errorf("customerService -> CreateCustomer - CustomerRepository.Create: %w", err)
	}

	return nil, err
}
