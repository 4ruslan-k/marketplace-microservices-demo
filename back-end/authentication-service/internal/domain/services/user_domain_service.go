package domainservices

import (
	userEntity "authentication_service/internal/domain/entities/user"
	"authentication_service/internal/ports/repositories"
	"context"
	"fmt"

	customErrors "authentication_service/pkg/errors"

	"github.com/rs/zerolog"
)

var ErrorEmailIsTaken = customErrors.NewIncorrectInputError("email_is_taken", "This Email is taken by another account")

var _ UserDomainService = (*userService)(nil)

type UserDomainService interface {
	CreateUser(ctx context.Context, params userEntity.CreateUserParams) (*userEntity.User, error)
}

type userService struct {
	logger                      zerolog.Logger
	authenticationDomainService AuthenticationDomainService
	userRepository              repositories.UserRepository
}

func NewUserService(
	logger zerolog.Logger,
	authenticationDomainService AuthenticationDomainService,
	userRepository repositories.UserRepository,
) *userService {
	return &userService{logger, authenticationDomainService, userRepository}
}

// Creates a user from user input or from social account
func (u userService) CreateUser(ctx context.Context, createUserParams userEntity.CreateUserParams) (*userEntity.User, error) {
	user, err := u.userRepository.GetByEmail(ctx, createUserParams.Email)
	if err != nil {
		return nil, fmt.Errorf("userService -> CreateUser-GetByEmail: %w", err)
	}
	if user != nil {
		return nil, ErrorEmailIsTaken
	}
	newUser, err := u.createUserEntityWithHashedPassword(createUserParams)
	if err != nil {
		return nil, fmt.Errorf("userService -> CreateUser-createUser: %w", err)
	}

	userID, err := u.userRepository.Create(ctx, *newUser)
	if err != nil {
		return nil, fmt.Errorf("userService -> CreateUser-userRepository.Create: %w", err)
	}
	createdUser, err := u.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("userService -> CreateUser-userRepository.GetByID: %w", err)
	}
	return createdUser, nil
}

func (u userService) createUserEntityWithHashedPassword(params userEntity.CreateUserParams) (*userEntity.User, error) {
	hashedPassword, err := u.authenticationDomainService.GetPasswordHashValue(params.Password)
	if err != nil {
		return nil, fmt.Errorf("userService -> createUserEntityWithHashedPassword - GetPasswordHashValue: %w", err)
	}
	newUser, err := userEntity.NewUser(userEntity.CreateUserParams{
		Name:          params.Name,
		Email:         params.Email,
		Password:      hashedPassword,
		SocialAccount: params.SocialAccount,
	})
	if err != nil {
		return nil, fmt.Errorf("userService -> createUserEntityWithHashedPassword - NewUser: %w", err)
	}
	return newUser, nil
}
