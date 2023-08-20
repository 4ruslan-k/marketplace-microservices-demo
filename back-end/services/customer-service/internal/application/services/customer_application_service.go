package applicationservices

import (
	"context"
	customerEntity "customer_service/internal/domain/entities/customer"
	domainServices "customer_service/internal/domain/services"
	"customer_service/internal/ports/repositories"
	"fmt"
	"time"

	"github.com/rs/zerolog"

	domainDto "customer_service/internal/application/dto"
	customErrors "shared/errors"
)

var ErrNoCustomer = customErrors.NewIncorrectInputError("no_customer", "No customer with this email")
var ErrInvalidCredentials = customErrors.NewIncorrectInputError("invalid_credentials", "Invalid credentials")
var ErrCustomerNotCreated = customErrors.NewIncorrectInputError("customer_not_created", "User wasn't created")
var ErrCustomerNotFound = customErrors.NewIncorrectInputError("no_customer", "User not found")

var _ CustomerApplicationService = (*customerApplicationService)(nil)

type customerApplicationService struct {
	customerRepository    repositories.CustomerRepository
	customerDomainService domainServices.CustomerDomainService
	logger                zerolog.Logger
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
	customerDomainService domainServices.CustomerDomainService,
) CustomerApplicationService {
	return customerApplicationService{customerRepository, customerDomainService, logger}
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
	createdCustomer, err := u.customerDomainService.CreateCustomer(ctx, createCustomerParams)
	if err != nil {
		return nil, err
	}
	return CustomerEntityToOutput(createdCustomer), nil
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
	fmt.Printf("update custoemr %+v\n", customer)
	fmt.Printf("update custoemr %+v\n", customer)
	fmt.Printf("update custoemr %+v\n", customer)
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
