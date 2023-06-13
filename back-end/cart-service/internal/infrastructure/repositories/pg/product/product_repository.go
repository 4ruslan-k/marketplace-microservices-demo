package pgrepositories

import (
	productEntity "cart_service/internal/domain/entities/product"
	"cart_service/internal/ports/repositories"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

type ProductModel struct {
	bun.BaseModel `bun:"table:products"`

	ID        string    `bun:"id,pk"`
	Price     float64   `bun:"price"`
	Name      string    `bun:"name"`
	Quantity  int       `bun:"quantity"`
	CreatedAt time.Time `bun:"created_at,nullzero"`
	UpdatedAt time.Time `bun:"updated_at,nullzero"`
}

var _ repositories.ProductRepository = (*productPGRepository)(nil)

type productPGRepository struct {
	db     *bun.DB
	logger zerolog.Logger
}

func (p *ProductModel) toDB(product productEntity.Product) ProductModel {
	return ProductModel{
		ID:        product.ID(),
		Name:      product.Name(),
		Price:     product.Price(),
		Quantity:  product.Quantity(),
		CreatedAt: product.CreatedAt(),
		UpdatedAt: product.UpdatedAt(),
	}
}

func (p *ProductModel) toEntity() productEntity.Product {

	product := productEntity.NewProductFromDatabase(
		p.ID,
		p.Price,
		p.Quantity,
		p.Name,
		p.CreatedAt,
		p.UpdatedAt,
	)

	return product
}

func NewProductRepository(sql *bun.DB, logger zerolog.Logger) *productPGRepository {
	return &productPGRepository{sql, logger}
}

func (r *productPGRepository) CreateProduct(ctx context.Context, product productEntity.Product) error {
	var dbProduct ProductModel

	dbProduct = dbProduct.toDB(product)
	_, err := r.db.NewInsert().Model(&dbProduct).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *productPGRepository) GetProductByID(ctx context.Context, id string) (productEntity.Product, error) {

	var productDB ProductModel
	err := r.db.NewSelect().
		Model(&productDB).
		Where("id IN (?)", id).
		Scan(ctx)

	if err == sql.ErrNoRows {
		return productEntity.Product{}, nil
	}

	if err != nil {
		return productEntity.Product{}, fmt.Errorf("customerPGRepository -> GetByID -> r.db.NewSelect(): %w", err)
	}

	product := productDB.toEntity()
	if err != nil {
		return productEntity.Product{}, fmt.Errorf("customerPGRepository -> GetByID -> toEntity: %w", err)
	}
	return product, nil
}
