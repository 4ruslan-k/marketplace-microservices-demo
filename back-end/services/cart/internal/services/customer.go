package applicationservices

import (
	customerEntity "cart/internal/domain/entities/customer"
	repositories "cart/internal/repositories/customer"
	repository "cart/internal/repositories/customer"
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"

	domainDto "cart/internal/services/dto"
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
	customerRepository repositories.CustomerRepository,
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

func (u customerApplicationService) createCustomer(
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
