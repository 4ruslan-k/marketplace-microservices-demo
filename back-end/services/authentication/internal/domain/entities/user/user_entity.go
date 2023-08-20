package user

import (
	"fmt"
	"regexp"
	"time"

	customErrors "shared/errors"

	mfaSettingsEntity "authentication/internal/domain/entities/mfa_settings"
	socialAccountEntity "authentication/internal/domain/entities/social_account"

	"github.com/google/uuid"
)

var ErrInvalidPassword = customErrors.NewIncorrectInputError("invalid_password", "invalid password")
var ErrInvalidEmailFormat = customErrors.NewIncorrectInputError("invalid_email", "invalid email format")

type User struct {
	id             string
	name           string
	email          string
	password       string
	mfaSettings    mfaSettingsEntity.MfaSettings
	socialAccounts []socialAccountEntity.SocialAccount
	createdAt      time.Time
	updatedAt      *time.Time
}

type CreateUserParams struct {
	Name          string
	Email         string
	Password      string
	SocialAccount *socialAccountEntity.CreateSocialAccountParams
}

func ValidateEmail(email string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	isValid := re.MatchString(email)
	return isValid
}

func NewUser(createUserParams CreateUserParams) (*User, error) {
	id := uuid.New().String()

	isValidEmail := ValidateEmail(createUserParams.Email)

	if !isValidEmail {
		return nil, ErrInvalidEmailFormat
	}

	if createUserParams.Password == "" {
		return nil, ErrInvalidPassword
	}

	var socialAccounts []socialAccountEntity.SocialAccount
	if createUserParams.SocialAccount != nil {
		var socialAccount *socialAccountEntity.SocialAccount
		socialAccount, err := socialAccountEntity.NewSocialAccount(
			*createUserParams.SocialAccount,
		)
		if err != nil {
			return nil, fmt.Errorf("NewUser -> socialAccountEntity.NewSocialAccount %w", err)
		}
		socialAccounts = append(socialAccounts, *socialAccount)
	}

	mfaSettings := mfaSettingsEntity.NewMfaSettings(false, "")

	createdAt := time.Now()
	user := User{id: id,
		email:          createUserParams.Email,
		name:           createUserParams.Name,
		password:       createUserParams.Password,
		createdAt:      createdAt,
		updatedAt:      nil,
		mfaSettings:    mfaSettings,
		socialAccounts: socialAccounts,
	}
	return &user, nil
}

func NewUserFromDatabase(
	id string,
	email string,
	name string,
	password string,
	createdAt time.Time,
	UpdatedAt *time.Time,
	socialAccounts []socialAccountEntity.SocialAccount,
	mfaSettings mfaSettingsEntity.MfaSettings,
) (*User, error) {
	user := User{id: id,
		email:          email,
		name:           name,
		password:       password,
		mfaSettings:    mfaSettings,
		createdAt:      createdAt,
		updatedAt:      nil,
		socialAccounts: socialAccounts,
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

func (u User) Password() string {
	return u.password
}

func (u User) CreatedAt() time.Time {
	return u.createdAt
}

func (u User) UpdatedAt() *time.Time {
	return u.updatedAt
}

func (u User) SocialAccounts() []socialAccountEntity.SocialAccount {
	return u.socialAccounts
}

func (u *User) SetName(name string) {
	u.name = name
}

func (u *User) SetPasswordHash(password string) {
	u.password = password
}

func (u *User) MfaSettings() *mfaSettingsEntity.MfaSettings {
	return &u.mfaSettings
}
