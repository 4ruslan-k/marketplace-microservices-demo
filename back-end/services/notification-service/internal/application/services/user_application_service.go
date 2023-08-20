package applicationservices

import (
	"context"
	"fmt"
	userEntity "notification_service/internal/domain/entities/user"
	domainServices "notification_service/internal/domain/services"
	"notification_service/internal/ports/repositories"
	"time"

	"github.com/rs/zerolog"

	domainDto "notification_service/internal/application/dto"
	customErrors "notification_service/pkg/errors"
)

var ErrNoUser = customErrors.NewIncorrectInputError("no_user", "No user with this email")
var ErrInvalidCredentials = customErrors.NewIncorrectInputError("invalid_credentials", "Invalid credentials")
var ErrUserNotCreated = customErrors.NewIncorrectInputError("user_not_created", "User wasn't created")
var ErrUserNotFound = customErrors.NewIncorrectInputError("no_user", "User not found")

var _ UserApplicationService = (*userApplicationService)(nil)

type userApplicationService struct {
	userRepository    repositories.UserRepository
	userDomainService domainServices.UserDomainService
	logger            zerolog.Logger
}

func UserEntityToOutput(user *userEntity.User) *domainDto.UserOutput {
	if user == nil {
		return nil
	}
	return &domainDto.UserOutput{
		ID:    user.ID(),
		Name:  user.Name(),
		Email: user.Email(),
	}
}

type UserApplicationService interface {
	GetUserByID(ctx context.Context, ID string) (*domainDto.UserOutput, error)
	CreateUser(ctx context.Context, createUser userEntity.CreateUserParams) (*domainDto.UserOutput, error)
	UpdateUser(ctx context.Context, updateUser domainDto.UpdateUserInput) error
	DeleteUser(ctx context.Context, deleteUser domainDto.DeleteUserInput) error
}

func NewUserApplicationService(
	userRepository repositories.UserRepository,
	logger zerolog.Logger,
	userDomainService domainServices.UserDomainService,
) userApplicationService {
	return userApplicationService{userRepository, userDomainService, logger}
}

func (u userApplicationService) GetUserByID(ctx context.Context, userID string) (*domainDto.UserOutput, error) {
	user, err := u.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return UserEntityToOutput(user), err
}

// Creates a user from user input
func (u userApplicationService) CreateUser(
	ctx context.Context,
	createUserParams userEntity.CreateUserParams,
) (*domainDto.UserOutput, error) {
	createdUser, err := u.userDomainService.CreateUser(ctx, createUserParams)
	if err != nil {
		return nil, err
	}
	return UserEntityToOutput(createdUser), nil
}

// Updates a user
func (u userApplicationService) UpdateUser(
	ctx context.Context,
	updateUserInput domainDto.UpdateUserInput,
) error {
	user, err := u.userRepository.GetByID(ctx, updateUserInput.ID)
	if err != nil {
		return fmt.Errorf("userApplicationService -> UpdateUser - u.userRepository.GetByID: %w", err)
	}

	if user == nil {
		return ErrUserNotFound
	}

	if len(updateUserInput.Name) > 0 {
		user.SetName(updateUserInput.Name)
	}

	user.SetUpdatedAt(time.Now())

	err = u.userRepository.Update(ctx, *user)
	if err != nil {
		return fmt.Errorf("userApplicationService -> UpdateUser - u.userRepository.Update: %w", err)
	}

	return nil
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

	return nil
}
