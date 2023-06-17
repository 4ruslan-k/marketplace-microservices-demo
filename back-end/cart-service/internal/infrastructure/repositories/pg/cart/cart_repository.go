package pgrepositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	cartEntity "cart_service/internal/domain/entities/cart"
	productRepo "cart_service/internal/infrastructure/repositories/pg/product"
	"cart_service/internal/ports/repositories"

	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

var _ repositories.CartRepository = (*cartRepository)(nil)

type CartModel struct {
	CustomerID string
	Products   []CartProductModel
}

type CartProductModel struct {
	bun.BaseModel `bun:"table:cart_products,alias:cart_products"`
	CustomerID    string                   `bun:"customer_id,pk"`
	ProductID     string                   `bun:"product_id,pk"`
	Quantity      int                      `bun:"quantity"`
	CreatedAt     time.Time                `bun:"created_at"`
	UpdatedAt     time.Time                `bun:"updated_at"`
	Product       productRepo.ProductModel `bun:"rel:has-one,join:product_id=id"`
}

type cartRepository struct {
	logger zerolog.Logger
	db     *bun.DB
}

func NewCartRepository(sql *bun.DB, logger zerolog.Logger) *cartRepository {
	return &cartRepository{logger: logger, db: sql}
}

func (r *cartRepository) GetByCustomerID(ctx context.Context, customerID string) (cartEntity.CartReadModel, error) {

	var cartProducts []CartProductModel
	err := r.db.NewSelect().
		Model(&cartProducts).
		ColumnExpr("cart_products.*").
		Join("JOIN products AS p ON p.id = cart_products.product_id").
		Where("cart_products.customer_id = ?", customerID).
		Scan(ctx)

	if err != nil {
		if err != sql.ErrNoRows {
			return cartEntity.CartReadModel{}, fmt.Errorf("cartProductRepository -> GetByCustomerID -> r.db.NewSelect(): %w", err)
		}
	}

	cartReadModelProducts := make([]cartEntity.CartReadModelProduct, 0, len(cartProducts))
	for _, cartProduct := range cartProducts {
		cartReadModelProduct := cartEntity.CartReadModelProduct{
			ProductID: cartProduct.ProductID,
			Quantity:  cartProduct.Quantity,
			Name:      cartProduct.Product.Name,
			Price:     cartProduct.Product.Price,
		}
		cartReadModelProducts = append(cartReadModelProducts, cartReadModelProduct)
	}

	cart := cartEntity.NewCartReadModel(
		customerID,
		cartReadModelProducts,
	)
	return cart, nil
}
