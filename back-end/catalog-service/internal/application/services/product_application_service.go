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
	GetProducts(ctx context.Context) ([]productEntity.Product, error)
	GetProductByID(ctx context.Context, productID string) (productEntity.Product, error)
	DeleteProductByID(ctx context.Context, productID string) error
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

	err = u.productRepository.SaveProduct(ctx, product)
	if err != nil {
		return fmt.Errorf("productApplicationService -> CreateProduct -  u.productRepository.Sav: %w", err)
	}

	return nil
}

// Fetches a list of products
func (u productApplicationService) GetProducts(
	ctx context.Context,
) ([]productEntity.Product, error) {
	products, err := u.productRepository.GetProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("productApplicationService -> GetProducts - u.productRepository.GetProducts: %w", err)
	}

	return products, nil
}

// Fetches products
func (u productApplicationService) GetProductByID(
	ctx context.Context,
	productID string,
) (productEntity.Product, error) {
	product, err := u.productRepository.GetProductByID(ctx, productID)
	if err != nil {
		return productEntity.Product{}, fmt.Errorf("productApplicationService -> GetProducts - u.productRepository.GetProducts: %w", err)
	}

	return product, nil
}

// Fetches products
func (u productApplicationService) DeleteProductByID(
	ctx context.Context,
	productID string,
) error {
	err := u.productRepository.DeleteProductByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("productApplicationService -> GetProducts - u.productRepository.DeleteProductByID: %w", err)
	}

	return nil
}
