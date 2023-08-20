package domainservices

import (
	userEntity "analytics_service/internal/domain/entities/user"
	repositories "analytics_service/internal/repositories/user"
	"context"
	"fmt"

	"github.com/rs/zerolog"
)

type UserDomainService interface {
	CreateUser(ctx context.Context, params userEntity.CreateUserParams) (*userEntity.User, error)
}

type userService struct {
	logger         zerolog.Logger
	userRepository repositories.UserRepository
}

func NewUserService(
	logger zerolog.Logger,
	userRepository repositories.UserRepository,
) *userService {
	return &userService{logger, userRepository}
}

// Creates a user from user input or from social account
func (u userService) CreateUser(ctx context.Context, createUserParams userEntity.CreateUserParams) (*userEntity.User, error) {
	u.logger.Info().Interface("createUserParams", createUserParams).Msg("userService -> CreateUser -> createUserParams")
	newUser, err := userEntity.NewUser(userEntity.CreateUserParams{
		Name:  createUserParams.Name,
		Email: createUserParams.Email,
		ID:    createUserParams.ID,
	})

	if err != nil {
		return nil, fmt.Errorf("userService -> CreateUser - NewUser: %w", err)
	}

	err = u.userRepository.Create(ctx, *newUser)

	if err != nil {
		return nil, fmt.Errorf("userService -> CreateUser - UserRepository.Create: %w", err)
	}

	return nil, err
}
