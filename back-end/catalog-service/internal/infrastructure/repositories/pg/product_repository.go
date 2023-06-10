package pgrepositories

import (
	productEntity "catalog_service/internal/domain/entities/product"
	"catalog_service/internal/ports/repositories"
	"context"
	"database/sql"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

var _ repositories.ProductRepository = (*productPGRepository)(nil)

type ProductModel struct {
	bun.BaseModel `bun:"table:products,alias:p"`

	ID       string  `bun:"id"`
	Name     string  `bun:"name"`
	Price    float64 `bun:"price"`
	Quantity int     `bun:"quantity"`
}

type productPGRepository struct {
	db     *bun.DB
	logger zerolog.Logger
}

func (p ProductModel) toEntity() productEntity.Product {

	product := productEntity.NewProductFromDatabase(
		p.ID,
		p.Name,
		p.Price,
		p.Quantity,
	)

	return product
}

func toDB(p productEntity.Product) (ProductModel, error) {
	return ProductModel{
		ID:    p.ID(),
		Name:  p.Name(),
		Price: p.Price(),
	}, nil
}

func NewProductRepository(sql *bun.DB, logger zerolog.Logger) *productPGRepository {
	return &productPGRepository{sql, logger}
}

func (r *productPGRepository) SaveProduct(ctx context.Context, u productEntity.Product) error {
	dbProduct, err := toDB(u)
	if err != nil {
		return err
	}
	_, err = r.db.NewInsert().Model(&dbProduct).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *productPGRepository) GetProducts(ctx context.Context) ([]productEntity.Product, error) {
	productModels := make([]ProductModel, 0)
	err := r.db.NewSelect().
		Model(&productModels).
		Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("productRepo -> GetProducts -> r.db.NewSelect(): %w", err)
	}

	products := make([]productEntity.Product, 0, len(productModels))
	for _, productModel := range productModels {
		productEntity := productModel.toEntity()

		products = append(products, productEntity)
	}

	return products, nil
}

func (r *productPGRepository) GetProductByID(ctx context.Context, productID string) (productEntity.Product, error) {
	var productDB ProductModel
	err := r.db.NewSelect().
		Model(&productDB).
		Where("id IN (?)", productID).
		Scan(ctx)

	if err == sql.ErrNoRows {
		return productEntity.Product{}, nil
	}

	if err != nil {
		return productEntity.Product{}, fmt.Errorf("productPGRepository -> GetProductByID -> r.db.NewSelect(): %w", err)
	}

	product := productDB.toEntity()
	if err != nil {
		return product, fmt.Errorf("productPGRepository -> GetByID -> toEntity: %w", err)
	}
	return product, nil
}

func (r *productPGRepository) DeleteProductByID(ctx context.Context, productID string) error {
	var product ProductModel
	_, err := r.db.NewDelete().Model(&product).Where("id = ?", productID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("productPGRepository DeleteProductByID -> NewDelete: %w", err)
	}

	return nil
}
