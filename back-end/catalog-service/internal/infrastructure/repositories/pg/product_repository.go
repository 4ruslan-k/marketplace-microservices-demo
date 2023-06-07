package pgrepositories

import (
	productEntity "catalog_service/internal/domain/entities/product"
	"catalog_service/internal/ports/repositories"
	"context"

	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

var _ repositories.ProductRepository = (*productPGRepository)(nil)

type ProductModel struct {
	bun.BaseModel `bun:"table:products,alias:p"`

	ID   string `bun:"id"`
	Name string `bun:"name"`
}

type productPGRepository struct {
	db     *bun.DB
	logger zerolog.Logger
}

func toEntity(p ProductModel) (productEntity.Product, error) {

	product, err := productEntity.NewProductFromDatabase(
		p.ID,
		p.Name,
	)
	if err != nil {
		return productEntity.Product{}, err
	}

	return product, nil
}

func toDB(u productEntity.Product) (ProductModel, error) {
	return ProductModel{
		ID:   u.ID(),
		Name: u.Name(),
	}, nil
}

func NewProductRepository(sql *bun.DB, logger zerolog.Logger) *productPGRepository {
	return &productPGRepository{sql, logger}
}

func (r *productPGRepository) Save(ctx context.Context, u productEntity.Product) error {
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
