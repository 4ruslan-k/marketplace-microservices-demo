package soccialaccount

import (
	"regexp"
	customErrors "shared/errors"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidEmailFormat = customErrors.NewIncorrectInputError("", "invalid email format")
var ErrInvalidProvider = customErrors.NewIncorrectInputError("", "invalid provider format")

type SocialAccount struct {
	id        string
	name      string
	email     string
	provider  string
	createdAt time.Time
	updatedAt *time.Time
}

type CreateSocialAccountParams struct {
	ID        string
	Name      string
	Email     string
	Provider  string
	createdAt time.Time
	updatedAt *time.Time
}

func NewSocialAccountFromDatabase(
	id string,
	name string,
	email string,
	provider string,
	createdAt time.Time,
	updatedAt *time.Time,
) (*SocialAccount, error) {
	socialAccount := SocialAccount{
		id:        id,
		email:     email,
		name:      name,
		provider:  provider,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
	return &socialAccount, nil
}

func ValidateEmail(email string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	isValid := re.MatchString(email)
	return isValid
}

func NewSocialAccount(createSocialAccount CreateSocialAccountParams) (*SocialAccount, error) {
	id := uuid.New().String()

	isValidEmail := ValidateEmail(createSocialAccount.Email)

	if !isValidEmail {
		return nil, ErrInvalidEmailFormat
	}

	if createSocialAccount.Provider == "" {
		return nil, ErrInvalidProvider
	}

	socialAccount := SocialAccount{
		id:        id,
		email:     createSocialAccount.Email,
		name:      createSocialAccount.Name,
		provider:  createSocialAccount.Provider,
		createdAt: time.Now(),
		updatedAt: nil,
	}
	return &socialAccount, nil
}

func (u SocialAccount) ID() string {
	return u.id
}

func (u SocialAccount) Name() string {
	return u.name
}

func (u SocialAccount) Email() string {
	return u.email
}

func (u SocialAccount) Provider() string {
	return u.provider
}
func (u SocialAccount) CreatedAt() time.Time {
	return u.createdAt
}

func (u SocialAccount) UpdatedAt() *time.Time {
	return u.updatedAt
}
