package applicationservices

import (
	"context"
	"fmt"

	productEntity "cart_service/internal/domain/entities/product"
	"cart_service/internal/ports/repositories"
	nats "cart_service/pkg/messaging/nats"

	"github.com/rs/zerolog"
)

var _ ProductApplicationService = (*productApplicationService)(nil)

type productApplicationService struct {
	productRepository repositories.ProductRepository
	logger            zerolog.Logger
	natsClient        nats.NatsClient
}

type ProductApplicationService interface {
	CreateProduct(ctx context.Context, createProductParams productEntity.CreateProductParams) error
	AddProductToCart(ctx context.Context) error
	DeleteProduct(ctx context.Context) error
}

func NewProductApplicationService(
	productRepository repositories.ProductRepository,
	logger zerolog.Logger,
	natsClient nats.NatsClient,
) productApplicationService {
	return productApplicationService{productRepository, logger, natsClient}
}

func (n productApplicationService) CreateProduct(
	ctx context.Context,
	createProductParams productEntity.CreateProductParams,
) error {
	product, err := productEntity.NewProduct(createProductParams)
	if err != nil {
		return fmt.Errorf("productApplicationService CreateProduct ->  productEntity.NewProduct: %w", err)
	}
	err = n.productRepository.CreateProduct(ctx, product)
	if err != nil {
		return fmt.Errorf("productApplicationService CreateProduct -> productRepository.CreateProduct: %w", err)
	}
	return nil
}

func (n productApplicationService) AddProductToCart(
	ctx context.Context,
) error {
	return nil
}

func (n productApplicationService) DeleteProduct(ctx context.Context) error {

	return nil
}
