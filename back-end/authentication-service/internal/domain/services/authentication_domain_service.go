package domainservices

import (
	passwordVerificationEntity "authentication_service/internal/domain/entities/password_verification_token"
	"authentication_service/internal/ports/repositories"
	customErrors "authentication_service/pkg/errors"
	"bytes"
	"context"
	"fmt"
	"image/png"
	"time"

	"encoding/base64"

	"github.com/pquerna/otp/totp"
	"github.com/rs/zerolog"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidPassword = customErrors.NewIncorrectInputError("invalid_password", "invalid password")

var _ AuthenticationDomainService = (*authenticationDomainService)(nil)

type AuthenticationDomainService interface {
	GetPasswordHashValue(password string) (string, error)
	VerifyPassword(userPassword string, providedPassword string) error
	GenerateTotp(email string) (TotpSetupInfo, error)
	ValidateTotp(otp string, otpSecretKey string) bool
	GenerateAndSavePasswordVerificationToken(ctx context.Context, userID string) (string, error)
}

type authenticationDomainService struct {
	authenticationRepository repositories.AuthenticationRepository
}

func GetPasswordHashValue(password string) (string, error) {
	if password == "" {
		return "", ErrInvalidPassword
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	hashedString := string(hash)
	return hashedString, err
}

func NewAuthenticationService(logger zerolog.Logger, authenticationRepository repositories.AuthenticationRepository) *authenticationDomainService {
	return &authenticationDomainService{authenticationRepository}
}

func (a authenticationDomainService) GetPasswordHashValue(password string) (string, error) {
	return GetPasswordHashValue(password)
}

func (a authenticationDomainService) VerifyPassword(userPassword string, providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}

type TotpSetupInfo struct {
	Image  string
	Secret string
}

func (a authenticationDomainService) GenerateTotp(email string) (TotpSetupInfo, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Marketplace Microservices Demo",
		AccountName: email,
	})
	if err != nil {
		return TotpSetupInfo{}, fmt.Errorf("authenticationDomainService -> totp.Generate: %w", err)
	}

	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		return TotpSetupInfo{}, fmt.Errorf("authenticationDomainService -> key.Image: %w", err)
	}
	png.Encode(&buf, img)
	imgBase64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
	otp := TotpSetupInfo{
		Image:  imgBase64Str,
		Secret: key.Secret(),
	}

	return otp, nil
}

// validate otp
func (a authenticationDomainService) ValidateTotp(otp string, otpSecretKey string) bool {
	valid := totp.Validate(otp, otpSecretKey)
	return valid
}

// Generates and saves password verification token for user
// Password verification tokens are used to skip password verification once the password was verified
func (a authenticationDomainService) GenerateAndSavePasswordVerificationToken(ctx context.Context, userID string) (string, error) {
	passwordTokenID, err := password.Generate(120, 10, 10, false, true)
	if err != nil {
		return "", fmt.Errorf("authenticationDomainService -> password.Generate: %w", err)
	}
	passwordVerificationToken := passwordVerificationEntity.NewPasswordVerificationToken(
		passwordVerificationEntity.CreatePasswordVerificationToken{
			ID:                 passwordTokenID,
			UserID:             userID,
			CurrentTime:        time.Now(),
			ExpirationDuration: passwordVerificationEntity.TokenExpirationDuration,
		},
	)
	err = a.authenticationRepository.SavePasswordVerificationToken(ctx, passwordVerificationToken)
	if err != nil {
		return "", fmt.Errorf("authenticationDomainService -> generateAndSavePasswordVerificationToken a.authenticationRepository.SavePasswordVerificationToken: %w", err)
	}
	return passwordTokenID, nil
}
