package product

import (
	customErrors "catalog_service/pkg/errors"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidProductName = customErrors.NewIncorrectInputError("products/invalid_name", "Invalid product name")
var ErrInvalidProductPrice = customErrors.NewIncorrectInputError("products/invalid_price", "Invalid product name")

type Product struct {
	id        string
	name      string
	price     float64
	quantity  int
	createdAt time.Time
	updatedAt time.Time
}

type CreateProductParams struct {
	Name     string
	Price    float64
	Quantity int
}

func NewProduct(createProductParams CreateProductParams) (Product, error) {
	id := uuid.New().String()

	if createProductParams.Name == "" {
		return Product{}, ErrInvalidProductName
	}

	if createProductParams.Price == 0 {
		return Product{}, ErrInvalidProductPrice
	}

	product := Product{
		id:        id,
		name:      createProductParams.Name,
		price:     createProductParams.Price,
		createdAt: time.Now(),
	}

	return product, nil
}

func NewProductFromDatabase(id, name string, price float64, quantity int, createdAt time.Time, updatedAt time.Time) Product {
	product := Product{
		id:        id,
		name:      name,
		price:     price,
		quantity:  quantity,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}

	return product
}

func (p Product) ID() string {
	return p.id
}

func (p Product) Name() string {
	return p.name
}

func (p Product) Price() float64 {
	return p.price
}

func (p Product) Quantity() int {
	return p.quantity
}

func (p Product) CreatedAt() time.Time {
	return p.createdAt
}
func (p Product) UpdatedAt() time.Time {
	return p.updatedAt
}

func (p Product) IsZero() bool {
	return p == Product{}
}
