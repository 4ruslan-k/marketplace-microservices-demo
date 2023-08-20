package applicationservices

import (
	socialAccountEntity "authentication/internal/domain/entities/social_account"
	userEntity "authentication/internal/domain/entities/user"
	domainServices "authentication/internal/domain/services"
	authRepository "authentication/internal/repositories/authentication"
	userRepository "authentication/internal/repositories/user"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	domainDto "authentication/internal/services/dto"
	customErrors "shared/errors"
	nats "shared/messaging/nats"
)

var (
	ErrNoUserByEmail                = customErrors.NewIncorrectInputError("no_user_by_email", "No user with this email")
	ErrNoUserByID                   = customErrors.NewIncorrectInputError("no_user_by_id", "User not found")
	ErrUserNotFound                 = customErrors.NewIncorrectInputError("no_user", "User not found")
	ErrInvalidCredentials           = customErrors.NewIncorrectInputError("invalid_credentials", "Invalid credentials")
	ErrUserNotCreated               = customErrors.NewIncorrectInputError("user_not_created", "User wasn't created")
	ErrPasswordsDoNotMatch          = customErrors.NewIncorrectInputError("passwords_do_not_match", "Passwords do not match")
	ErrTotpCodeNotValid             = customErrors.NewIncorrectInputError("totp_not_valid", "Not valid code")
	ErrTotpMfaAlreadyActiveNotValid = customErrors.NewIncorrectInputError("totp_already_enabled", "TOTP MFA already enabled")
	ErrTotpMfaNotEnabled            = customErrors.NewIncorrectInputError("totp_not_enabled", "TOTP MFA is not enabled")
)

var _ UserApplicationService = (*userApplicationService)(nil)

type userApplicationService struct {
	userRepository              userRepository.UserRepository
	authenticationRepository    authRepository.AuthenticationRepository
	userDomainService           domainServices.UserDomainService
	logger                      zerolog.Logger
	authenticationDomainService domainServices.AuthenticationDomainService
	natsClient                  nats.NatsClient
}

func UserEntityToOutput(user *userEntity.User) *domainDto.UserOutput {
	if user == nil {
		return nil
	}
	return &domainDto.UserOutput{
		ID:           user.ID(),
		Name:         user.Name(),
		Email:        user.Email(),
		IsMfaEnabled: user.MfaSettings().IsMfaEnabled(),
	}
}

type UserApplicationService interface {
	GetUserByID(ctx context.Context, ID string) (*domainDto.UserOutput, error)
	CreateUser(ctx context.Context, createUser userEntity.CreateUserParams) (*domainDto.UserOutput, error)
	UpdateUser(ctx context.Context, updateUser domainDto.UpdateUserInput) (*domainDto.UserOutput, error)
	DeleteUser(ctx context.Context, deleteUser domainDto.DeleteUserInput) error
	LoginWithEmailAndPassword(ctx context.Context, email string, password string) (domainDto.LoginOutput, error)
	LoginWithTotpCode(ctx context.Context, passwordVerificationTokenID string, code string) (domainDto.UserOutput, error)
	SocialLogin(ctx context.Context, socialAccount domainDto.SocialLoginInput) (*domainDto.LoginOutput, error)
	ChangeCurrentPassword(ctx context.Context, socialAccount domainDto.ChangeCurrentPasswordInput) error
	GenerateTotpSetup(ctx context.Context, userID string) (domainServices.TotpSetupInfo, error)
	EnableTotp(ctx context.Context, userID string, otp string) error
	DisableTotp(ctx context.Context, userID string, otp string) error
}

func NewUserApplicationService(
	userRepository userRepository.UserRepository,
	authenticationRepository authRepository.AuthenticationRepository,
	logger zerolog.Logger,
	userDomainService domainServices.UserDomainService,
	authenticationDomainService domainServices.AuthenticationDomainService,
	natsClient nats.NatsClient,
) userApplicationService {
	return userApplicationService{userRepository, authenticationRepository, userDomainService, logger, authenticationDomainService, natsClient}
}

func (u userApplicationService) GetUserByID(ctx context.Context, userID string) (*domainDto.UserOutput, error) {
	user, err := u.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return UserEntityToOutput(user), err
}

type UserCreatedEvent struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	ID    string `json:"id"`
}

type UserUpdatedEvent struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type UserDeletedEvent struct {
	ID string `json:"id"`
}

const (
	userNotificationCreationSubject = "notifications.created"
)

var (
	MFAEnabledNotificationTypeID  = "mfa-enabled-v1"
	MFADisabledNotificationTypeID = "mfa-disabled-v1"
)

type NotificationCreatedEvent struct {
	NotificationTypeID string      `json:"notificationTypeID"`
	UserID             string      `json:"userID"`
	Data               interface{} `json:"data,omitempty"`
}

// Creates a user from user input
func (u userApplicationService) CreateUser(
	ctx context.Context,
	createUserParams userEntity.CreateUserParams,
) (*domainDto.UserOutput, error) {
	createdUser, err := u.userDomainService.CreateUser(ctx, createUserParams)
	if err != nil {
		return nil, fmt.Errorf("userApplicationService -> CreateUser -  u.userDomainService.CreateUser: %w", err)
	}

	bytes, err := json.Marshal(UserCreatedEvent{
		Name:  createdUser.Name(),
		Email: createdUser.Email(),
		ID:    createdUser.ID(),
	})

	if err != nil {
		return nil, fmt.Errorf("userApplicationService -> CreateUser -  json.Marshal: %w", err)
	}

	u.natsClient.PublishMessage("users.created", string(bytes))
	return UserEntityToOutput(createdUser), nil
}

// Updates a user
func (u userApplicationService) UpdateUser(
	ctx context.Context,
	updateUserInput domainDto.UpdateUserInput,
) (*domainDto.UserOutput, error) {
	user, err := u.userRepository.GetByID(ctx, updateUserInput.ID)
	if err != nil {
		return nil, fmt.Errorf("userApplicationService -> UpdateUser - u.userRepository.GetByID: %w", err)
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	if len(updateUserInput.Name) > 0 {
		user.SetName(updateUserInput.Name)
	}

	err = u.userRepository.Update(ctx, *user)
	if err != nil {
		return nil, fmt.Errorf("userApplicationService -> UpdateUser - u.userRepository.Update: %w", err)
	}

	updatedUser, err := u.userRepository.GetByID(ctx, updateUserInput.ID)
	if err != nil {
		return nil, fmt.Errorf("userApplicationService -> UpdateUser - u.userRepository.GetByID: %w", err)
	}
	bytes, err := json.Marshal(UserUpdatedEvent{
		Name: updatedUser.Name(),
		ID:   updatedUser.ID(),
	})
	if err != nil {
		return nil, err
	}
	u.natsClient.PublishMessage("users.updated", string(bytes))
	return UserEntityToOutput(updatedUser), nil
}

// Deletes a user
func (u userApplicationService) DeleteUser(
	ctx context.Context,
	deleteUserInput domainDto.DeleteUserInput,
) error {
	user, err := u.userRepository.GetByID(ctx, deleteUserInput.ID)
	if err != nil {
		return fmt.Errorf("userApplicationService -> DeleteUser - u.userRepository.GetByID: %w", err)
	}

	if user == nil {
		return ErrUserNotFound
	}

	err = u.userRepository.Delete(ctx, deleteUserInput.ID)
	if err != nil {
		return fmt.Errorf("userApplicationService -> DeleteUser - u.userRepository.Delete: %w", err)
	}

	bytes, err := json.Marshal(UserDeletedEvent{
		ID: deleteUserInput.ID,
	})
	if err != nil {
		return fmt.Errorf("userApplicationService -> DeleteUser - json.Marshal: %w", err)
	}
	u.natsClient.PublishMessage("users.deleted", string(bytes))
	return nil
}

func (u userApplicationService) LoginWithEmailAndPassword(ctx context.Context, email string, password string) (domainDto.LoginOutput, error) {
	user, err := u.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return domainDto.LoginOutput{}, fmt.Errorf("userApplicationService -> LoginWithEmailAndPassword - GetByEmail: %w", err)
	}
	if user == nil {
		return domainDto.LoginOutput{}, ErrNoUserByEmail

	}
	err = u.authenticationDomainService.VerifyPassword(user.Password(), password)
	if err != nil {
		return domainDto.LoginOutput{}, ErrInvalidCredentials
	}

	var passwordVerificationTokenID string
	if user.MfaSettings().IsMfaEnabled() {
		passwordVerificationTokenID, err = u.authenticationDomainService.GenerateAndSavePasswordVerificationToken(ctx, user.ID())
		if err != nil {
			return domainDto.LoginOutput{}, fmt.Errorf("userApplicationService -> LoginWithEmailAndPassword - GenerateAndSavePasswordVerificationToken: %w", err)
		}
		return domainDto.LoginOutput{
			IsMfaEnabled:                user.MfaSettings().IsMfaEnabled(),
			PasswordVerificationTokenID: passwordVerificationTokenID,
		}, nil
	}

	return domainDto.LoginOutput{
		ID:                          user.ID(),
		Email:                       user.Email(),
		Name:                        user.Name(),
		IsMfaEnabled:                user.MfaSettings().IsMfaEnabled(),
		PasswordVerificationTokenID: passwordVerificationTokenID,
	}, nil
}

func (u userApplicationService) LoginWithTotpCode(ctx context.Context, passwordVerificationTokenID string, code string) (domainDto.UserOutput, error) {
	passwordVerificationToken, err := u.authenticationRepository.GetPasswordVerificationTokenByID(ctx, passwordVerificationTokenID)
	if err != nil {
		return domainDto.UserOutput{}, fmt.Errorf("userApplicationService -> LoginWithTotpCode - GetPasswordVerificationTokenByID: %w", err)
	}

	if passwordVerificationToken.IsZero() {
		return domainDto.UserOutput{}, ErrTotpCodeNotValid
	}

	isExpired := passwordVerificationToken.HasExpired(time.Now())
	if isExpired {
		return domainDto.UserOutput{}, ErrTotpCodeNotValid
	}

	user, err := u.userRepository.GetByID(ctx, passwordVerificationToken.UserID())
	if err != nil {
		return domainDto.UserOutput{}, fmt.Errorf("userApplicationService -> LoginWithTotpCode - GetByID: %w", err)
	}

	isValid := u.authenticationDomainService.ValidateTotp(code, user.MfaSettings().TotpSecret())
	if !isValid {
		return domainDto.UserOutput{}, ErrTotpCodeNotValid
	}
	err = u.authenticationRepository.DeletePasswordVerificationTokenByID(ctx, passwordVerificationTokenID)
	if err != nil {
		return domainDto.UserOutput{}, fmt.Errorf("userApplicationService -> LoginWithTotpCode - DeletePasswordVerificationTokenByID: %w", err)
	}

	return domainDto.UserOutput{
		ID:    user.ID(),
		Email: user.Email(),
		Name:  user.Name(),
	}, nil
}

func (u userApplicationService) SocialLogin(
	ctx context.Context,
	socialAccount domainDto.SocialLoginInput,
) (*domainDto.LoginOutput, error) {
	user, err := u.userRepository.GetByEmail(ctx, socialAccount.Email)
	if err != nil {
		return nil, fmt.Errorf("userApplicationService -> SocialLogin - u.userRepository.GetByEmail: %w", err)
	}
	var passwordVerificationTokenID string
	if user == nil {
		_, err := u.userDomainService.CreateUser(ctx, userEntity.CreateUserParams{
			Name:     socialAccount.Name,
			Email:    socialAccount.Email,
			Password: uuid.New().String() + uuid.New().String(),
			// TODO: use NewSocialAccountSettings
			SocialAccount: &socialAccountEntity.CreateSocialAccountParams{
				ID:       socialAccount.UserID,
				Name:     socialAccount.Name,
				Email:    socialAccount.Email,
				Provider: socialAccount.Provider,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("userApplicationService -> SocialLogin - u.userDomainService.CreateUser: %w", err)
		}

		user, err = u.userRepository.GetByEmail(ctx, socialAccount.Email)
		if user == nil {
			return nil, ErrUserNotCreated
		}
		if err != nil {
			return nil, fmt.Errorf("userApplicationService -> SocialLogin - u.userRepository.GetByEmail: %w", err)
		}
	} else {
		if user.MfaSettings().IsMfaEnabled() {
			passwordVerificationTokenID, err = u.authenticationDomainService.GenerateAndSavePasswordVerificationToken(ctx, user.ID())
			return &domainDto.LoginOutput{
				IsMfaEnabled:                user.MfaSettings().IsMfaEnabled(),
				PasswordVerificationTokenID: passwordVerificationTokenID,
			}, nil
		}
		if err != nil {
			return nil, fmt.Errorf("userApplicationService -> LoginWithEmailAndPassword - GenerateAndSavePasswordVerificationToken: %w", err)
		}

		// TODO: add social account
	}

	return &domainDto.LoginOutput{
		ID:                          user.ID(),
		Email:                       user.Email(),
		Name:                        user.Name(),
		IsMfaEnabled:                user.MfaSettings().IsMfaEnabled(),
		PasswordVerificationTokenID: passwordVerificationTokenID,
	}, nil
}

// Changes user's password with a new one after validating their current password
func (u userApplicationService) ChangeCurrentPassword(
	ctx context.Context,
	changeCurrentPasswordInput domainDto.ChangeCurrentPasswordInput,
) error {

	if changeCurrentPasswordInput.NewPassword != changeCurrentPasswordInput.NewPasswordConfirmation {
		return ErrPasswordsDoNotMatch
	}

	user, err := u.userRepository.GetByID(ctx, changeCurrentPasswordInput.UserID)
	if err != nil {
		return fmt.Errorf("userApplicationService -> ChangeCurrentPassword - u.userRepository.GetByID: %w", err)
	}

	if user == nil {
		return ErrNoUserByID

	}
	err = u.authenticationDomainService.VerifyPassword(user.Password(), changeCurrentPasswordInput.CurrentPassword)
	if err != nil {
		return ErrInvalidCredentials
	}

	newPasswordHash, err := u.authenticationDomainService.GetPasswordHashValue(changeCurrentPasswordInput.NewPassword)
	if err != nil {
		return fmt.Errorf("userApplicationService -> ChangeCurrentPassword - u.authenticationDomainService.GetPasswordHashValue: %w", err)
	}

	user.SetPasswordHash(newPasswordHash)
	err = u.userRepository.Update(ctx, *user)
	if err != nil {
		return fmt.Errorf("userApplicationService -> ChangeCurrentPassword - u.userRepository.Update: %w", err)
	}

	return nil
}

// Generates "Time-based one-time password" secret and image, saves secret in the DB
func (u userApplicationService) GenerateTotpSetup(
	ctx context.Context,
	userID string,
) (domainServices.TotpSetupInfo, error) {
	user, err := u.userRepository.GetByID(ctx, userID)
	if err != nil {
		return domainServices.TotpSetupInfo{}, fmt.Errorf("userApplicationService -> GenerateTotpSetup -u.userRepository.GetByID: %w", err)
	}

	if user == nil {
		return domainServices.TotpSetupInfo{}, ErrNoUserByID
	}

	IsMfaEnabled := user.MfaSettings().IsMfaEnabled()
	if IsMfaEnabled {
		return domainServices.TotpSetupInfo{}, ErrTotpMfaAlreadyActiveNotValid
	}

	totpInfo, err := u.authenticationDomainService.GenerateTotp(user.Email())
	if err != nil {
		return domainServices.TotpSetupInfo{}, fmt.Errorf("userApplicationService -> GenerateTotpSetup -u.authenticationDomainService.GenerateTotp: %w", err)
	}

	user.MfaSettings().SetTotpSecret(totpInfo.Secret)
	err = u.userRepository.Update(ctx, *user)
	if err != nil {
		return domainServices.TotpSetupInfo{}, fmt.Errorf("userApplicationService -> GenerateTotpSetup -u.userRepository.Update: %w", err)
	}
	return totpInfo, nil
}

// Enables MFA by validating the provided OTP
func (u userApplicationService) EnableTotp(
	ctx context.Context,
	userID string,
	otp string,
) error {

	user, err := u.userRepository.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("userApplicationService -> EnableTotp - u.userRepository.GetByID: %w", err)
	}

	if user == nil {
		return ErrNoUserByID
	}

	IsMfaEnabled := user.MfaSettings().IsMfaEnabled()
	if IsMfaEnabled {
		return ErrTotpMfaAlreadyActiveNotValid
	}

	isValid := u.authenticationDomainService.ValidateTotp(otp, user.MfaSettings().TotpSecret())

	if !isValid {
		return ErrTotpCodeNotValid
	}

	user.MfaSettings().SetMfaStatus(true)
	err = u.userRepository.Update(ctx, *user)
	if err != nil {
		return fmt.Errorf("userApplicationService -> EnableTotp - u.userRepository.Update: %w", err)
	}

	bytes, err := json.Marshal(NotificationCreatedEvent{
		UserID:             user.ID(),
		NotificationTypeID: MFAEnabledNotificationTypeID,
	})

	u.natsClient.PublishMessage(userNotificationCreationSubject, string(bytes))

	return nil
}

// Enables MFA by validating the provided OTP
func (u userApplicationService) DisableTotp(
	ctx context.Context,
	userID string,
	otp string,
) error {

	user, err := u.userRepository.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("userApplicationService -> DisableTotp - u.userRepository.GetByID: %w", err)
	}

	if user == nil {
		return ErrNoUserByID
	}

	IsMfaEnabled := user.MfaSettings().IsMfaEnabled()
	if !IsMfaEnabled {
		return ErrTotpMfaNotEnabled
	}

	isValid := u.authenticationDomainService.ValidateTotp(otp, user.MfaSettings().TotpSecret())

	if !isValid {
		return ErrTotpCodeNotValid
	}

	user.MfaSettings().SetMfaStatus(false)
	user.MfaSettings().SetTotpSecret("")
	err = u.userRepository.Update(ctx, *user)
	if err != nil {
		return fmt.Errorf("userApplicationService -> EnableTotp - u.userRepository.Update: %w", err)
	}

	bytes, err := json.Marshal(NotificationCreatedEvent{
		UserID:             user.ID(),
		NotificationTypeID: MFADisabledNotificationTypeID,
	})

	u.natsClient.PublishMessage(userNotificationCreationSubject, string(bytes))

	return nil
}
