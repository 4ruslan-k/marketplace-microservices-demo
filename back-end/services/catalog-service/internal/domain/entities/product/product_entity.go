package product

import (
	customErrors "catalog_service/pkg/errors"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidProductName = customErrors.NewIncorrectInputError("products/invalid_name", "Invalid product name")
var ErrInvalidProductPrice = customErrors.NewIncorrectInputError("products/invalid_price", "Invalid product name")
var ErrInvalidProductQuantity = customErrors.NewIncorrectInputError("products/invalid_quantity", "Invalid product quantity")

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

type UpdateProductParams struct {
	Name     *string
	Price    *float64
	Quantity *int
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
		quantity:  createProductParams.Quantity,
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

func (p Product) Update(params UpdateProductParams) (Product, error) {
	var isUpdated bool
	if params.Name != nil {
		if *params.Name == "" {
			return Product{}, ErrInvalidProductName
		}
		p.name = *params.Name
		isUpdated = true

	}

	if params.Price != nil {
		if *params.Price <= 0 {
			return Product{}, ErrInvalidProductPrice
		}
		p.price = *params.Price
		isUpdated = true

	}

	if params.Quantity != nil {
		if *params.Quantity < 0 {
			return Product{}, ErrInvalidProductQuantity
		}
		p.quantity = *params.Quantity
		isUpdated = true

	}

	if isUpdated {
		p.updatedAt = time.Now()
	}

	return p, nil
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
