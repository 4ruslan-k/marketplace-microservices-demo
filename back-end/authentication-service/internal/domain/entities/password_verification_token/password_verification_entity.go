package passwordveificationtoken

import (
	"time"
)

var TokenExpirationDuration = 10 * time.Minute

type PasswordVerificationToken struct {
	id        string
	userID    string
	createdAt time.Time
	expiresAt time.Time
}

func (p PasswordVerificationToken) ID() string {
	return p.id
}

func (p PasswordVerificationToken) UserID() string {
	return p.userID
}

func (p PasswordVerificationToken) CreatedAt() time.Time {
	return p.createdAt
}

func (p PasswordVerificationToken) ExpiresAt() time.Time {
	return p.expiresAt
}

func (p PasswordVerificationToken) HasExpired(currentTime time.Time) bool {
	return currentTime.After(p.expiresAt)
}

func (p PasswordVerificationToken) IsZero() bool {
	return p == PasswordVerificationToken{}
}

type CreatePasswordVerificationToken struct {
	ID                 string
	UserID             string
	CurrentTime        time.Time
	ExpirationDuration time.Duration
}

func NewPasswordVerificationTokeFromDatabase(
	id string,
	userID string,
	createdAt time.Time,
	expiresAt time.Time,
) PasswordVerificationToken {
	passwordVerificationToken := PasswordVerificationToken{
		id:        id,
		userID:    userID,
		createdAt: createdAt,
		expiresAt: expiresAt,
	}
	return passwordVerificationToken
}

func NewPasswordVerificationToken(createPasswordVerificationToken CreatePasswordVerificationToken) PasswordVerificationToken {
	passwordVerificationToken := PasswordVerificationToken{
		id:        createPasswordVerificationToken.ID,
		userID:    createPasswordVerificationToken.UserID,
		createdAt: createPasswordVerificationToken.CurrentTime,
		expiresAt: createPasswordVerificationToken.CurrentTime.Add(createPasswordVerificationToken.ExpirationDuration),
	}
	return passwordVerificationToken
}
