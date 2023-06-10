package user

import (
	"regexp"
	"time"

	customErrors "notification_service/pkg/errors"
)

var ErrInvalidEmailFormat = customErrors.NewIncorrectInputError("invalid_email", "invalid email format")

type User struct {
	id        string
	name      string
	email     string
	password  string
	createdAt time.Time
	updatedAt time.Time
}

type CreateUserParams struct {
	Name  string
	Email string
	ID    string
}

func ValidateEmail(email string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	isValid := re.MatchString(email)
	return isValid
}

func NewUser(createUserParams CreateUserParams) (*User, error) {

	user := User{id: createUserParams.ID,
		email:     createUserParams.Email,
		name:      createUserParams.Name,
		createdAt: time.Now(),
	}
	return &user, nil
}

func NewUserFromDatabase(
	id string,
	email string,
	name string,
) (*User, error) {
	user := User{id: id,
		email: email,
		name:  name,
	}
	return &user, nil
}

func (u User) ID() string {
	return u.id
}

func (u User) Name() string {
	return u.name
}

func (u User) Email() string {
	return u.email
}

func (u User) CreatedAt() time.Time {
	return u.createdAt
}

func (u User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) SetName(name string) {
	u.name = name
}

func (u *User) SetUpdatedAt(updatedAt time.Time) {
	u.updatedAt = updatedAt
}
