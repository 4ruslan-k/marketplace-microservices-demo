package product

import (
	"time"
)

type Product struct {
	id        string
	name      string
	price     float64
	quantity  int
	createdAt time.Time
	updatedAt time.Time
}

type CreateProductParams struct {
	ID        string
	Name      string
	Price     float64
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewProduct(createProductParams CreateProductParams) (Product, error) {
	product := Product{
		id:        createProductParams.ID,
		name:      createProductParams.Name,
		price:     createProductParams.Price,
		quantity:  createProductParams.Quantity,
		createdAt: createProductParams.CreatedAt,
		updatedAt: createProductParams.UpdatedAt,
	}
	return product, nil
}

func NewProductFromDatabase(
	id string,
	price float64,
	quantity int,
	name string,
	createdAt time.Time,
	updatedAt time.Time,
) Product {
	product := Product{
		id:        id,
		name:      name,
		quantity:  quantity,
		price:     price,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
	return product
}

func (d Product) ID() string {
	return d.id
}

func (d Product) Price() float64 {
	return d.price
}

func (d Product) Quantity() int {
	return d.quantity
}

func (d Product) Name() string {
	return d.name
}

func (d Product) CreatedAt() time.Time {
	return d.createdAt
}

func (d Product) UpdatedAt() time.Time {
	return d.updatedAt
}

func (d Product) IsZero() bool {
	return d == Product{}
}
