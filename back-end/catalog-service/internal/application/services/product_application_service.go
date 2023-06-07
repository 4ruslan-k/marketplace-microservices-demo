package applicationservices

import (
	productEntity "catalog_service/internal/domain/entities/product"
	"catalog_service/internal/ports/repositories"
	"context"
	"fmt"

	"github.com/rs/zerolog"
)

var _ ProductApplicationService = (*productApplicationService)(nil)

type productApplicationService struct {
	productRepository repositories.ProductRepository
	logger            zerolog.Logger
}

type ProductApplicationService interface {
	CreateProduct(ctx context.Context, createProductParams productEntity.CreateProductParams) error
}

func NewProductApplicationService(
	productRepository repositories.ProductRepository,
	logger zerolog.Logger,
) *productApplicationService {
	return &productApplicationService{productRepository: productRepository, logger: logger}
}

// Creates a product
func (u productApplicationService) CreateProduct(
	ctx context.Context,
	createProductParams productEntity.CreateProductParams,
) error {
	product, err := productEntity.NewProduct(createProductParams)
	if err != nil {
		return fmt.Errorf("productApplicationService -> CreateProduct - productEntity.NewProduct: %w", err)
	}

	err = u.productRepository.Save(ctx, product)
	if err != nil {
		return fmt.Errorf("productApplicationService -> CreateProduct -  u.productRepository.Sav: %w", err)
	}

	return nil
}
