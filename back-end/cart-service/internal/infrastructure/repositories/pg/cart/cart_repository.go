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

func (r *cartRepository) SaveCart(
	ctx context.Context,
	customerID string,
	updateFunc func(cart cartEntity.Cart,
	) (cartEntity.Cart, error)) error {
	var cartProductModels []CartProductModel
	err := r.db.NewSelect().
		Model(&cartProductModels).
		ColumnExpr("cart_products.*").
		Join("JOIN products AS p ON p.id = cart_products.product_id").
		Where("cart_products.customer_id = ?", customerID).
		Scan(ctx)

	if err != nil {
		if err != sql.ErrNoRows {
			return fmt.Errorf("cartProductRepository -> SaveCart -> r.db.NewSelect(): %w", err)
		}
	}

	cartProducts := make([]cartEntity.CartProduct, 0, len(cartProductModels))
	for _, cartProduct := range cartProductModels {
		product := cartEntity.CartProduct{
			ProductID: cartProduct.ProductID,
			Quantity:  cartProduct.Quantity,
		}
		cartProducts = append(cartProducts, product)
	}

	cart, err := cartEntity.NewCart(cartEntity.CreateCartParams{
		CustomerID: customerID,
		Products:   cartProducts,
	})
	if err != nil {
		return fmt.Errorf("cartRepository -> SaveCart -> cartEntity.NewCart(): %w", err)
	}
	cart, err = updateFunc(cart)
	if err != nil {
		return fmt.Errorf("cartRepository -> SaveCart -> updateFunc(): %w", err)
	}

	events := cart.Events()

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("cartRepository -> SaveCart -> r.db.BeginTx(): %w", err)
	}

	for _, event := range events {
		if ev, ok := event.(cartEntity.ProductAdded); ok {
			product := CartProductModel{
				CustomerID: cart.CustomerID(),
				ProductID:  ev.Product.ProductID,
				Quantity:   ev.Product.Quantity,
			}
			_, err := r.db.NewInsert().Model(&product).Exec(ctx)
			if err != nil {
				txError := tx.Rollback()
				if txError != nil {
					return fmt.Errorf("cartRepository -> SaveCart -> tx.Rollback(): %w", txError)
				}
				return fmt.Errorf("cartRepository -> SaveCart -> r.db.NewInsert(): %w", err)
			}
		} else if ev, ok := event.(cartEntity.ProductRemoved); ok {
			_, err := r.db.NewDelete().
				Model(&CartProductModel{}).
				Where("customer_id = ? AND product_id = ?", cart.CustomerID(), ev.ProductID).
				Exec(ctx)
			if err != nil {
				txError := tx.Rollback()
				if txError != nil {
					return fmt.Errorf("cartRepository -> SaveCart -> tx.Rollback(): %w", txError)
				}
				return fmt.Errorf("cartRepository -> SaveCart ->r.db.NewDelete(): %w", err)
			}
		} else if ev, ok := event.(cartEntity.ProductQuantityChanged); ok {
			_, err := r.db.NewUpdate().
				Model(&CartProductModel{}).
				Set("quantity = ?", ev.Product.Quantity).
				Where("customer_id = ? AND product_id = ?", cart.CustomerID(), ev.Product.ProductID).
				Exec(ctx)
			if err != nil {
				txError := tx.Rollback()
				if txError != nil {
					return fmt.Errorf("cartRepository -> SaveCart -> tx.Rollback(): %w", txError)
				}
				return fmt.Errorf("cartRepository -> SaveCart -> r.db.NewUpdate(): %w", err)
			}
		} else {
			return fmt.Errorf("cartRepository -> SaveCart -> unknown event: %v", event)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("cartRepository -> SaveCart -> tx.Commit(): %w", err)
	}
	return nil
}
