package repository

import (
	"context"
	"database/sql"
	"fmt"
	userEntity "notification_service/internal/domain/entities/user"
	repository "notification_service/internal/repositories/user"
	"time"

	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

type UserModel struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        string    `bun:"id"`
	Name      string    `bun:"name"`
	Email     string    `bun:"email"`
	CreatedAt time.Time `bun:"created_at,nullzero"`
	UpdatedAt time.Time `bun:"updated_at,nullzero"`
}

var _ repository.UserRepository = (*userPGRepository)(nil)

type userPGRepository struct {
	db     *bun.DB
	logger zerolog.Logger
}

func toEntity(u UserModel) (*userEntity.User, error) {

	user, err := userEntity.NewUserFromDatabase(
		u.ID,
		u.Email,
		u.Name,
		u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func toDB(u userEntity.User) (UserModel, error) {
	return UserModel{
		ID:        u.ID(),
		Name:      u.Name(),
		Email:     u.Email(),
		CreatedAt: u.CreatedAt(),
		UpdatedAt: u.UpdatedAt(),
	}, nil
}

func NewUserRepository(sql *bun.DB, logger zerolog.Logger) *userPGRepository {
	return &userPGRepository{sql, logger}
}

func (r *userPGRepository) GetByID(ctx context.Context, id string) (*userEntity.User, error) {

	var user UserModel
	err := r.db.NewSelect().
		Model(&user).
		Where("id IN (?)", id).
		Scan(ctx)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("userPGRepository -> GetByID -> r.db.NewSelect(): %w", err)
	}

	userEntity, err := toEntity(user)
	if err != nil {
		return nil, fmt.Errorf("userPGRepository -> GetByID -> toEntity: %w", err)
	}
	return userEntity, nil
}

func (r *userPGRepository) GetByEmail(ctx context.Context, email string) (*userEntity.User, error) {
	var user UserModel
	err := r.db.NewSelect().
		Model(&user).
		Where("email IN (?)", email).
		Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("userPGRepository -> GetByEmail -> r.db.NewSelect(): %w", err)
	}
	userEntity, err := toEntity(user)
	if err != nil {
		return nil, fmt.Errorf("userPGRepository -> GetByEmail -> toEntity(user): %w", err)
	}
	return userEntity, nil
}

func (r *userPGRepository) Create(ctx context.Context, u userEntity.User) error {
	dbUser, err := toDB(u)
	if err != nil {
		return err
	}
	_, err = r.db.NewInsert().Model(&dbUser).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *userPGRepository) Update(ctx context.Context, user userEntity.User) error {
	sqlUser, err := toDB(user)
	if err != nil {
		return fmt.Errorf("userPGRepository Update -> toDB(user): %w", err)
	}

	_, err = r.db.NewUpdate().Model(&sqlUser).Where("id = ?", sqlUser.ID).Exec(ctx)

	if err != nil {
		return fmt.Errorf("userPGRepository Update -> r.db.NewUpdate: %w", err)
	}

	return nil
}

func (r *userPGRepository) Delete(ctx context.Context, ID string) error {
	var user UserModel
	_, err := r.db.NewDelete().Model(&user).Where("id = ?", ID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("userPGRepository Delete -> NewDelete: %w", err)
	}

	return nil
}
