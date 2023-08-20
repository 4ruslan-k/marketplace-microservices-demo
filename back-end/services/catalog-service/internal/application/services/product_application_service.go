package applicationservices

import (
	productEntity "catalog_service/internal/domain/entities/product"
	"catalog_service/internal/ports/repositories"
	natsClient "catalog_service/pkg/nats"
	"context"
	"encoding/json"
	"fmt"
	"time"

	customErrors "catalog_service/pkg/errors"

	"github.com/rs/zerolog"
)

var _ ProductApplicationService = (*productApplicationService)(nil)

type productApplicationService struct {
	productRepository repositories.ProductRepository
	logger            zerolog.Logger
	natsClient        natsClient.NatsClient
}

type ProductCreatedEvent struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductDeletedEvent struct {
	ID string `json:"id"`
}

type ProductUpdatedEvent struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const (
	productCreatedSubject             = "products.created"
	productDeletedSubject             = "products.deleted"
	productUpdatedSubject             = "products.updated"
	productCreatedDurableConsumerName = "cart-service-product-created"
	productStream                     = "products"
)

type ProductApplicationService interface {
	CreateProduct(ctx context.Context, createProductParams productEntity.CreateProductParams) error
	GetProducts(ctx context.Context) ([]productEntity.Product, error)
	GetProductByID(ctx context.Context, productID string) (productEntity.Product, error)
	DeleteProductByID(ctx context.Context, productID string) error
	UpdateProductByID(ctx context.Context, productID string, productParams productEntity.UpdateProductParams) error
}

func NewProductApplicationService(
	productRepository repositories.ProductRepository,
	logger zerolog.Logger,
	nats natsClient.NatsClient,
) *productApplicationService {
	return &productApplicationService{productRepository: productRepository, logger: logger, natsClient: nats}
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

	bytes, err := json.Marshal(ProductCreatedEvent{
		ID:        product.ID(),
		Name:      product.Name(),
		Price:     product.Price(),
		Quantity:  product.Quantity(),
		CreatedAt: product.CreatedAt(),
		UpdatedAt: product.UpdatedAt(),
	})
	if err != nil {
		return err
	}
	u.natsClient.PublishMessage(productCreatedSubject, string(bytes))

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

// Fetches a product by ID
func (u productApplicationService) GetProductByID(
	ctx context.Context,
	productID string,
) (productEntity.Product, error) {
	product, err := u.productRepository.GetProductByID(ctx, productID)
	if product.IsZero() {
		return productEntity.Product{}, customErrors.NewNotFoundError("products/not_found", "product not found")
	}
	if err != nil {
		return productEntity.Product{}, fmt.Errorf("productApplicationService -> GetProducts - u.productRepository.GetProducts: %w", err)
	}

	return product, nil
}

// Deletes product by ID
func (u productApplicationService) DeleteProductByID(
	ctx context.Context,
	productID string,
) error {
	err := u.productRepository.DeleteProductByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("productApplicationService -> DeleteProductByID - u.productRepository.DeleteProductByID: %w", err)
	}
	bytes, err := json.Marshal(ProductDeletedEvent{
		ID: productID,
	})
	if err != nil {
		return err
	}
	u.natsClient.PublishMessage(productDeletedSubject, string(bytes))
	return nil
}

// Updates product by ID
func (u productApplicationService) UpdateProductByID(ctx context.Context, productID string, productParams productEntity.UpdateProductParams) error {
	product, err := u.productRepository.GetProductByID(ctx, productID)
	if product.IsZero() {
		return customErrors.NewNotFoundError("products/not_found", "product not found")
	}
	if err != nil {
		return fmt.Errorf("productApplicationService -> UpdateProductByID - u.productRepository.GetProductByID: %w", err)
	}
	product, err = product.Update(productParams)
	if err != nil {
		return fmt.Errorf("productApplicationService -> UpdateProductByID - product.Update: %w", err)
	}

	err = u.productRepository.UpdateProductByID(ctx, product)
	if err != nil {
		return fmt.Errorf("productApplicationService -> UpdateProductByID - u.productRepository.UpdateProductByID: %w", err)
	}
	bytes, err := json.Marshal(ProductUpdatedEvent{
		ID:        productID,
		Name:      product.Name(),
		Price:     product.Price(),
		Quantity:  product.Quantity(),
		CreatedAt: product.CreatedAt(),
		UpdatedAt: product.UpdatedAt(),
	})
	if err != nil {
		return err
	}
	u.natsClient.PublishMessage(productUpdatedSubject, string(bytes))
	return nil
}
