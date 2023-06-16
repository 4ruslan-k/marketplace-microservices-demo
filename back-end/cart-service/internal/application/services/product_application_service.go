package applicationservices

import (
	"context"
	"fmt"

	cartEntity "cart_service/internal/domain/entities/cart"
	productEntity "cart_service/internal/domain/entities/product"
	"cart_service/internal/ports/repositories"
	customErrors "cart_service/pkg/errors"
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
	GetCart(ctx context.Context, customerID string) (cartEntity.Cart, error)
	CreateProduct(ctx context.Context, createProductParams productEntity.CreateProductParams) error
	AddProductToCart(ctx context.Context, productID string, quantity int, customerID string) error
	DeleteProduct(ctx context.Context) error
}

func NewProductApplicationService(
	productRepository repositories.ProductRepository,
	logger zerolog.Logger,
	natsClient nats.NatsClient,
) productApplicationService {
	return productApplicationService{productRepository, logger, natsClient}
}

func (p productApplicationService) CreateProduct(
	ctx context.Context,
	createProductParams productEntity.CreateProductParams,
) error {
	product, err := productEntity.NewProduct(createProductParams)
	if err != nil {
		return fmt.Errorf("productApplicationService CreateProduct ->  productEntity.NewProduct: %w", err)
	}
	err = p.productRepository.CreateProduct(ctx, product)
	if err != nil {
		return fmt.Errorf("productApplicationService CreateProduct -> productRepository.CreateProduct: %w", err)
	}
	return nil
}

var ErrInvalidEmailFormat = customErrors.NewNotFoundError("cart/products", "Product not found")

func (p productApplicationService) AddProductToCart(
	ctx context.Context,
	productID string,
	quantity int,
	customerID string,
) error {
	product, err := p.productRepository.GetProductByID(ctx, productID)
	if product.IsZero() {
		return ErrInvalidEmailFormat
	}

	products := []cartEntity.CartProduct{
		{
			ProductID: productID,
			Quantity:  quantity,
		},
	}

	cart, err := cartEntity.NewCart(cartEntity.CreateCartParams{
		CustomerID: customerID,
		Products:   products,
	})

	if err != nil {
		return fmt.Errorf("productApplicationService AddProductToCart -> cartEntity.NewCart: %w", err)
	}

	cart, err = cart.AddProductToCart(
		cartEntity.CartProduct{
			ProductID: productID,
			Quantity:  quantity,
		},
	)

	if err != nil {
		return fmt.Errorf("productApplicationService AddProductToCart -> cartEntity.NewCart: %w", err)
	}

	// TODO: save card state using repository

	if err != nil {
		return fmt.Errorf("productApplicationService AddProductToCart -> productRepository.GetProductByID: %w", err)
	}

	return nil
}

func (p productApplicationService) GetCart(ctx context.Context, customerID string) (cartEntity.Cart, error) {
	panic("implement me")

	return cartEntity.Cart{}, nil
}

func (p productApplicationService) DeleteProduct(ctx context.Context) error {
	panic("implement me")
	return nil
}
