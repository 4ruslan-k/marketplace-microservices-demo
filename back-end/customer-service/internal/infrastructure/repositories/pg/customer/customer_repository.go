package pgrepositories

import (
	"context"
	customerEntity "customer_service/internal/domain/entities/customer"
	"customer_service/internal/ports/repositories"
	"database/sql"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

var _ repositories.CustomerRepository = (*customerPGRepository)(nil)

type CustomerModel struct {
	bun.BaseModel `bun:"table:customers,alias:c"`

	ID        string    `bun:"id"`
	Name      string    `bun:"name"`
	Email     string    `bun:"email"`
	CreatedAt time.Time `bun:"created_at,nullzero"`
	UpdatedAt time.Time `bun:"updated_at,nullzero"`
}

type customerPGRepository struct {
	db     *bun.DB
	logger zerolog.Logger
}

func toEntity(u CustomerModel) (*customerEntity.Customer, error) {

	customer, err := customerEntity.NewCustomerFromDatabase(
		u.ID,
		u.Email,
		u.Name,
		u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func toDB(u customerEntity.Customer) (CustomerModel, error) {
	return CustomerModel{
		ID:        u.ID(),
		Name:      u.Name(),
		Email:     u.Email(),
		CreatedAt: u.CreatedAt(),
		UpdatedAt: u.UpdatedAt(),
	}, nil
}

func NewCustomerRepository(sql *bun.DB, logger zerolog.Logger) *customerPGRepository {
	return &customerPGRepository{sql, logger}
}

func (r *customerPGRepository) GetByID(ctx context.Context, id string) (*customerEntity.Customer, error) {

	var customer CustomerModel
	err := r.db.NewSelect().
		Model(&customer).
		Where("id IN (?)", id).
		Scan(ctx)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("customerPGRepository -> GetByID -> r.db.NewSelect(): %w", err)
	}

	customerEntity, err := toEntity(customer)
	if err != nil {
		return nil, fmt.Errorf("customerPGRepository -> GetByID -> toEntity: %w", err)
	}
	return customerEntity, nil
}

func (r *customerPGRepository) GetByEmail(ctx context.Context, email string) (*customerEntity.Customer, error) {
	var customer CustomerModel
	// err := r.customersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&customer)
	err := r.db.NewSelect().
		Model(&customer).
		Where("email IN (?)", email).
		Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("customerPGRepository -> GetByEmail -> r.db.NewSelect(): %w", err)
	}
	customerEntity, err := toEntity(customer)
	if err != nil {
		return nil, fmt.Errorf("customerPGRepository -> GetByEmail -> toEntity(customer): %w", err)
	}
	return customerEntity, nil
}

func (r *customerPGRepository) Create(ctx context.Context, u customerEntity.Customer) error {
	dbUser, err := toDB(u)
	if err != nil {
		return err
	}
	_, err = r.db.NewInsert().Model(&dbUser).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *customerPGRepository) Update(ctx context.Context, customer customerEntity.Customer) error {
	sqlCustomer, err := toDB(customer)
	if err != nil {
		return fmt.Errorf("customerPGRepository Update -> toDB(customer): %w", err)
	}

	_, err = r.db.NewUpdate().Model(&sqlCustomer).Where("id = ?", sqlCustomer.ID).Exec(ctx)

	if err != nil {
		return fmt.Errorf("customerPGRepository Update -> r.db.NewUpdate: %w", err)
	}

	return nil
}

func (r *customerPGRepository) Delete(ctx context.Context, ID string) error {
	var customer CustomerModel
	_, err := r.db.NewDelete().Model(&customer).Where("id = ?", ID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("customerPGRepository Delete -> NewDelete: %w", err)
	}

	return nil
}
