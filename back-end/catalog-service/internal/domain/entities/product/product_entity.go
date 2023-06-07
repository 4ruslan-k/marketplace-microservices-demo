package product

import (
	customErrors "catalog_service/pkg/errors"

	"github.com/google/uuid"
)

var ErrInvalidProductName = customErrors.NewIncorrectInputError("products/invalid_name", "invalid name")

type Product struct {
	id   string
	name string
}

type CreateProductParams struct {
	Name string
}

func NewProduct(createProductParams CreateProductParams) (Product, error) {
	id := uuid.New().String()

	if createProductParams.Name == "" {
		return Product{}, ErrInvalidProductName
	}

	product := Product{
		id:   id,
		name: createProductParams.Name,
	}

	return product, nil
}

func NewProductFromDatabase(id, name string) (Product, error) {
	product := Product{
		id:   id,
		name: name,
	}

	return product, nil
}

func (p Product) ID() string {
	return p.id
}

func (p Product) Name() string {
	return p.name
}
