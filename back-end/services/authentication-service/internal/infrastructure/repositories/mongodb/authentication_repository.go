package mongorepositories

import (
	passwordVerificationTokenEntity "authentication_service/internal/domain/entities/password_verification_token"
	"authentication_service/internal/ports/repositories"
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ repositories.AuthenticationRepository = (*authenticationMongoDbRepository)(nil)

type PasswordVerificationTokenModel struct {
	ID        string             `bson:"_id,omitempty"`
	UserID    string             `bson:"userId,omitempty"`
	CreatedAt primitive.DateTime `bson:"createdAt,omitempty"`
	ExpiresAt primitive.DateTime `bson:"expiresAt,omitempty"`
}

type authenticationMongoDbRepository struct {
	mongoDB                              *mongo.Database
	passwordVerificationTokensCollection *mongo.Collection
	logger                               zerolog.Logger
}

func (p PasswordVerificationTokenModel) toEntity() passwordVerificationTokenEntity.PasswordVerificationToken {
	passwordVerificationToken := passwordVerificationTokenEntity.NewPasswordVerificationTokeFromDatabase(
		p.ID,
		p.UserID,
		p.CreatedAt.Time(),
		p.ExpiresAt.Time(),
	)
	return passwordVerificationToken
}

func (p PasswordVerificationTokenModel) fromEntity(pe passwordVerificationTokenEntity.PasswordVerificationToken) PasswordVerificationTokenModel {
	return PasswordVerificationTokenModel{
		ID:        pe.ID(),
		UserID:    pe.UserID(),
		CreatedAt: primitive.NewDateTimeFromTime(pe.CreatedAt()),
		ExpiresAt: primitive.NewDateTimeFromTime(pe.ExpiresAt()),
	}
}

func NewAuthenticationRepository(m *mongo.Database, logger zerolog.Logger) *authenticationMongoDbRepository {
	passwordVerificationTokensCollection := m.Collection("password_verification_tokens")
	return &authenticationMongoDbRepository{m, passwordVerificationTokensCollection, logger}
}

func (r *authenticationMongoDbRepository) SavePasswordVerificationToken(ctx context.Context, passwordVerificationToken passwordVerificationTokenEntity.PasswordVerificationToken) error {
	passwordVerificationTokenModel := PasswordVerificationTokenModel{}.fromEntity(passwordVerificationToken)
	_, err := r.passwordVerificationTokensCollection.InsertOne(ctx, passwordVerificationTokenModel)
	if err != nil {
		return fmt.Errorf("authenticationMongoDbRepository CreatePasswordVerificationToken -> InsertOne: %w", err)
	}
	return nil
}

func (r *authenticationMongoDbRepository) GetPasswordVerificationTokenByID(
	ctx context.Context,
	id string,
) (passwordVerificationTokenEntity.PasswordVerificationToken, error) {
	var passwordVerificationToken PasswordVerificationTokenModel
	err := r.passwordVerificationTokensCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&passwordVerificationToken)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return passwordVerificationTokenEntity.PasswordVerificationToken{}, nil
		}
		return passwordVerificationTokenEntity.PasswordVerificationToken{}, fmt.Errorf("authenticationMongoDbRepository GetPasswordVerificationTokenByID -> FindOne: %w", err)
	}
	passwordVerificationTokenEntity := passwordVerificationToken.toEntity()
	return passwordVerificationTokenEntity, nil
}

func (r *authenticationMongoDbRepository) DeletePasswordVerificationTokenByID(ctx context.Context, ID string) error {

	_, err := r.passwordVerificationTokensCollection.DeleteOne(ctx, bson.M{"_id": ID})
	if err != nil {
		return fmt.Errorf("authenticationMongoDbRepository DeletePasswordVerificationTokenByID -> DeleteOne: %w", err)
	}

	return nil
}
