package mongorepositories

import (
	mfaSettingsEntity "authentication_service/internal/domain/entities/mfa_settings"
	socialAccountEntity "authentication_service/internal/domain/entities/social_account"
	userEntity "authentication_service/internal/domain/entities/user"
	"authentication_service/internal/ports/repositories"
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserModel struct {
	ID             string               `bson:"_id,omitempty"`
	Name           string               `bson:"name,omitempty"`
	Email          string               `bson:"email,omitempty"`
	Password       string               `bson:"password,omitempty"`
	CreatedAt      primitive.DateTime   `bson:"createdAt,omitempty"`
	UpdatedAt      *primitive.DateTime  `bson:"updatedAt,omitempty"`
	SocialAccounts []SocialAccountModel `bson:"socialAccounts,omitempty"`
	MfaSettings    MfaSettingsModel     `bson:"mfaSettingsModel,omitempty"`
}

type SocialAccountModel struct {
	ID        string              `bson:"_id,omitempty"`
	Name      string              `bson:"name,omitempty"`
	Email     string              `bson:"email,omitempty"`
	Provider  string              `bson:"provider,omitempty"`
	CreatedAt primitive.DateTime  `bson:"createdAt,omitempty"`
	UpdatedAt *primitive.DateTime `bson:"updatedAt,omitempty"`
}

type MfaSettingsModel struct {
	IsMfaEnabled bool               `bson:"isMfaEnabled,omitempty"`
	TotpSecret   string             `bson:"totpSecret,omitempty"`
	CreatedAt    primitive.DateTime `bson:"createdAt,omitempty"`
	UpdatedAt    primitive.DateTime `bson:"updatedAt,omitempty"`
}

var _ repositories.UserRepository = (*userMongoDbRepository)(nil)

type userMongoDbRepository struct {
	mongoDB         *mongo.Database
	usersCollection *mongo.Collection
	logger          zerolog.Logger
}

func toEntity(u UserModel) (*userEntity.User, error) {
	socialAccounts := make([]socialAccountEntity.SocialAccount, len(u.SocialAccounts))
	for i, socialAccountMongo := range u.SocialAccounts {
		socialAccount, _ := socialAccountEntity.NewSocialAccountFromDatabase(
			socialAccountMongo.ID,
			socialAccountMongo.Name,
			socialAccountMongo.Email,
			socialAccountMongo.Provider,
			socialAccountMongo.CreatedAt.Time(),
			nil,
		)
		socialAccounts[i] = *socialAccount
	}

	mfaSettings := mfaSettingsEntity.NewMfaSettingsFromDatabase(
		u.MfaSettings.IsMfaEnabled,
		u.MfaSettings.TotpSecret,
		u.MfaSettings.CreatedAt.Time(),
		u.MfaSettings.UpdatedAt.Time(),
	)

	user, err := userEntity.NewUserFromDatabase(
		u.ID,
		u.Email,
		u.Name,
		u.Password,
		u.CreatedAt.Time(),
		nil,
		socialAccounts,
		mfaSettings,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func toMongoDB(u userEntity.User) (UserModel, error) {
	var socialAccountMongo []SocialAccountModel
	for _, v := range u.SocialAccounts() {
		socialAccountMongo = append(socialAccountMongo, SocialAccountModel{
			ID:        v.ID(),
			Name:      v.Name(),
			Email:     v.Email(),
			Provider:  v.Provider(),
			CreatedAt: primitive.NewDateTimeFromTime(v.CreatedAt()),
		})
	}
	mfaSettings := MfaSettingsModel{
		IsMfaEnabled: u.MfaSettings().IsMfaEnabled(),
		TotpSecret:   u.MfaSettings().TotpSecret(),
		CreatedAt:    primitive.NewDateTimeFromTime(u.MfaSettings().CreatedAt()),
		UpdatedAt:    primitive.NewDateTimeFromTime(u.MfaSettings().UpdatedAt()),
	}

	return UserModel{
		ID:        u.ID(),
		Name:      u.Name(),
		Email:     u.Email(),
		Password:  u.Password(),
		CreatedAt: primitive.NewDateTimeFromTime(u.CreatedAt()),
		// TODO: updateAt can be a value
		UpdatedAt:      nil,
		SocialAccounts: socialAccountMongo,
		MfaSettings:    mfaSettings,
	}, nil
}

func NewUserRepository(m *mongo.Database, logger zerolog.Logger) *userMongoDbRepository {
	usersCollection := m.Collection("users")
	return &userMongoDbRepository{m, usersCollection, logger}
}

func (r *userMongoDbRepository) GetByID(ctx context.Context, id string) (*userEntity.User, error) {
	var user UserModel
	err := r.usersCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("userMongoDbRepository GetByID -> GetByID: %w", err)
	}
	userEntity, err := toEntity(user)
	if err != nil {
		return nil, fmt.Errorf("userMongoDbRepository GetByID -> toEntity: %w", err)
	}
	return userEntity, nil
}

func (r *userMongoDbRepository) GetByEmail(ctx context.Context, email string) (*userEntity.User, error) {
	var user UserModel
	err := r.usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("userMongoDbRepository GetByEmail -> FindOne: %w", err)
	}
	userEntity, err := toEntity(user)
	if err != nil {
		return nil, fmt.Errorf("userMongoDbRepository GetByEmail -> toEntity: %w", err)
	}
	return userEntity, nil
}

func (r *userMongoDbRepository) Update(ctx context.Context, user userEntity.User) error {
	mongoUser, err := toMongoDB(user)
	if err != nil {
		return fmt.Errorf("userMongoDbRepository Update -> toMongoDB: %w", err)
	}

	_, err = r.usersCollection.ReplaceOne(ctx, bson.M{"_id": user.ID()}, mongoUser)
	if err != nil {
		return fmt.Errorf("userMongoDbRepository Update -> ReplaceOne: %w", err)
	}

	return nil
}

func (r *userMongoDbRepository) Delete(ctx context.Context, ID string) error {

	_, err := r.usersCollection.DeleteOne(ctx, bson.M{"_id": ID})
	if err != nil {
		return fmt.Errorf("userMongoDbRepository Delete -> DeleteOne: %w", err)
	}

	return nil
}

func (r *userMongoDbRepository) Create(ctx context.Context, u userEntity.User) (string, error) {
	mongoUser, err := toMongoDB(u)
	if err != nil {
		return "", fmt.Errorf("userMongoDbRepository Create -> toMongoDB: %w", err)
	}
	insertedRecord, err := r.usersCollection.InsertOne(ctx, mongoUser)
	if err != nil {
		return "", fmt.Errorf("userMongoDbRepository Create -> InsertOne: %w", err)
	}
	userID := insertedRecord.InsertedID.(string)
	return userID, nil
}
