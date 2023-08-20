package applicationservices

import (
	"context"
	customerEntity "customer/internal/domain/entities/customer"
	repository "customer/internal/repositories/customer"
	"fmt"
	"time"

	"github.com/rs/zerolog"

	domainDto "customer/internal/services/dto"
	customErrors "shared/errors"
)

var ErrNoCustomer = customErrors.NewIncorrectInputError("no_customer", "No customer with this email")
var ErrInvalidCredentials = customErrors.NewIncorrectInputError("invalid_credentials", "Invalid credentials")
var ErrCustomerNotCreated = customErrors.NewIncorrectInputError("customer_not_created", "User wasn't created")
var ErrCustomerNotFound = customErrors.NewIncorrectInputError("no_customer", "User not found")

var _ CustomerApplicationService = (*customerApplicationService)(nil)

type customerApplicationService struct {
	customerRepository repository.CustomerRepository
	logger             zerolog.Logger
}

func CustomerEntityToOutput(customer *customerEntity.Customer) *domainDto.CustomerOutput {
	if customer == nil {
		return nil
	}
	return &domainDto.CustomerOutput{
		ID:    customer.ID(),
		Name:  customer.Name(),
		Email: customer.Email(),
	}
}

type CustomerApplicationService interface {
	GetCustomerByID(ctx context.Context, ID string) (*domainDto.CustomerOutput, error)
	CreateCustomer(ctx context.Context, createCustomer customerEntity.CreateCustomerParams) (*domainDto.CustomerOutput, error)
	UpdateCustomer(ctx context.Context, updateCustomer domainDto.UpdateCustomerInput) error
	DeleteCustomer(ctx context.Context, deleteCustomer domainDto.DeleteCustomerInput) error
}

func NewCustomerApplicationService(
	customerRepository repository.CustomerRepository,
	logger zerolog.Logger,
) CustomerApplicationService {
	return customerApplicationService{customerRepository, logger}
}

func (u customerApplicationService) GetCustomerByID(ctx context.Context, customerID string) (*domainDto.CustomerOutput, error) {
	customer, err := u.customerRepository.GetByID(ctx, customerID)
	if err != nil {
		return nil, err
	}
	return CustomerEntityToOutput(customer), err
}

// Creates a customer from customer input
func (u customerApplicationService) CreateCustomer(
	ctx context.Context,
	createCustomerParams customerEntity.CreateCustomerParams,
) (*domainDto.CustomerOutput, error) {
	createdCustomer, err := u.createCustomer(ctx, createCustomerParams)
	if err != nil {
		return nil, err
	}
	return CustomerEntityToOutput(createdCustomer), nil
}

// Creates a customer from customer input or from social account
func (u customerApplicationService) createCustomer(ctx context.Context, createCustomerParams customerEntity.CreateCustomerParams) (*customerEntity.Customer, error) {
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

// Updates a customer
func (u customerApplicationService) UpdateCustomer(
	ctx context.Context,
	updateCustomerInput domainDto.UpdateCustomerInput,
) error {
	customer, err := u.customerRepository.GetByID(ctx, updateCustomerInput.ID)
	if err != nil {
		return fmt.Errorf("CustomerApplicationService -> UpdateCustomer - u.customerRepository.GetByID: %w", err)
	}

	if customer == nil {
		return ErrCustomerNotFound
	}

	if len(updateCustomerInput.Name) > 0 {
		customer.SetName(updateCustomerInput.Name)
	}
	customer.SetUpdatedAt(time.Now())
	err = u.customerRepository.Update(ctx, *customer)
	if err != nil {
		return fmt.Errorf("CustomerApplicationService -> UpdateCustomer - u.customerRepository.Update: %w", err)
	}

	return nil
}

// Deletes a customer
func (u customerApplicationService) DeleteCustomer(
	ctx context.Context,
	deleteCustomerInput domainDto.DeleteCustomerInput,
) error {
	customer, err := u.customerRepository.GetByID(ctx, deleteCustomerInput.ID)
	if err != nil {
		return fmt.Errorf("CustomerApplicationService -> DeleteCustomer - u.customerRepository.GetByID: %w", err)
	}

	if customer == nil {
		return ErrCustomerNotFound
	}

	err = u.customerRepository.Delete(ctx, deleteCustomerInput.ID)
	if err != nil {
		return fmt.Errorf("CustomerApplicationService -> DeleteCustomer - u.customerRepository.Delete: %w", err)
	}

	return nil
}
