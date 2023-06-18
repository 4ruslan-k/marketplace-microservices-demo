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
	cartRepository    repositories.CartRepository
	logger            zerolog.Logger
	natsClient        nats.NatsClient
}

type ProductApplicationService interface {
	GetCart(ctx context.Context, customerID string) (cartEntity.CartReadModel, error)
	CreateProduct(ctx context.Context, createProductParams productEntity.CreateProductParams) error
	DeleteProduct(ctx context.Context, productID string) error
	UpdateProductsInCart(ctx context.Context, productID string, quantity int, customerID string) error
}

func NewProductApplicationService(
	productRepository repositories.ProductRepository,
	cartRepository repositories.CartRepository,
	logger zerolog.Logger,
	natsClient nats.NatsClient,
) productApplicationService {
	return productApplicationService{
		productRepository: productRepository,
		cartRepository:    cartRepository,
		logger:            logger,
		natsClient:        natsClient,
	}
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

func (p productApplicationService) UpdateProductsInCart(
	ctx context.Context,
	productID string,
	quantity int,
	customerID string,
) error {
	updateOperation := func(cart cartEntity.Cart) (cartEntity.Cart, error) {

		cart, err := cart.UpdateProductsInCart(
			cartEntity.CartProduct{
				ProductID: productID,
				Quantity:  quantity,
			},
		)

		if err != nil {
			return cartEntity.Cart{}, fmt.Errorf("productApplicationService UpdateProductsInCart -> cartEntity.NewCart: %w", err)
		}
		return cart, nil
	}

	err := p.cartRepository.SaveCart(ctx, customerID, updateOperation)

	if err != nil {
		return fmt.Errorf("productApplicationService UpdateProductsInCart -> productRepository.GetProductByID: %w", err)
	}

	return nil
}

func (p productApplicationService) GetCart(ctx context.Context, customerID string) (cartEntity.CartReadModel, error) {
	cart, err := p.cartRepository.GetByCustomerID(ctx, customerID)
	if err != nil {
		return cartEntity.CartReadModel{}, fmt.Errorf("productApplicationService GetCart -> cartRepository.GetByCustomerID: %w", err)
	}
	return cart, nil
}

func (p productApplicationService) DeleteProduct(ctx context.Context, productID string) error {
	err := p.productRepository.DeleteProductByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("productApplicationService DeleteProduct -> p.productRepository.DeleteProductByID: %w", err)
	}
	return nil
}
