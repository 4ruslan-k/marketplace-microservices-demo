package Customer

import (
	"regexp"
	"time"

	customErrors "customer_service/pkg/errors"
)

var ErrInvalidEmailFormat = customErrors.NewIncorrectInputError("invalid_email", "invalid email format")

type Customer struct {
	id        string
	name      string
	email     string
	createdAt time.Time
	updatedAt time.Time
}

type CreateCustomerParams struct {
	Name  string
	Email string
	ID    string
}

func ValidateEmail(email string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	isValid := re.MatchString(email)
	return isValid
}

func NewCustomer(createCustomerParams CreateCustomerParams) (*Customer, error) {

	customer := Customer{id: createCustomerParams.ID,
		email:     createCustomerParams.Email,
		name:      createCustomerParams.Name,
		createdAt: time.Now(),
	}
	return &customer, nil
}

func NewCustomerFromDatabase(
	id string,
	email string,
	name string,
	createdAt time.Time,
) (*Customer, error) {
	customer := Customer{id: id,
		email:     email,
		name:      name,
		createdAt: createdAt,
	}
	return &customer, nil
}

func (u Customer) ID() string {
	return u.id
}

func (u Customer) Name() string {
	return u.name
}

func (u Customer) Email() string {
	return u.email
}

func (u Customer) CreatedAt() time.Time {
	return u.createdAt
}

func (u Customer) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *Customer) SetName(name string) {
	u.name = name
}

func (u *Customer) SetUpdatedAt(updatedAt time.Time) {
	u.updatedAt = updatedAt
}
