package repositories

import (
	passwordVerificationTokenEntity "authentication/internal/domain/entities/password_verification_token"
	"context"
)

type AuthenticationRepository interface {
	SavePasswordVerificationToken(ctx context.Context, passwordVerificationToken passwordVerificationTokenEntity.PasswordVerificationToken) error
	GetPasswordVerificationTokenByID(ctx context.Context, id string) (passwordVerificationTokenEntity.PasswordVerificationToken, error)
	DeletePasswordVerificationTokenByID(ctx context.Context, id string) error
}
