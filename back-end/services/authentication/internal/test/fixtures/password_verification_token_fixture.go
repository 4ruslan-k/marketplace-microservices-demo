package fixtures

import (
	passwordVerificationTokenEntity "authentication/internal/domain/entities/password_verification_token"
	"context"
	"testing"
	"time"
)

type CreatePasswordVerificationToken struct {
	ID        string
	UserID    string
	CreatedAt time.Time
	ExpiresAt time.Time
}

func GeneratePasswordVerificationTokenEntity(t *testing.T, c CreatePasswordVerificationToken) passwordVerificationTokenEntity.PasswordVerificationToken {
	t.Helper()
	id := c.ID
	userID := c.UserID
	createdAt := c.CreatedAt
	expiresAt := c.ExpiresAt

	if id == "" {
		id = fake.UUID().V4()
	}

	if userID == "" {
		id = fake.UUID().V4()
	}

	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	if expiresAt.IsZero() {
		expiresAt = createdAt.Add(1 * time.Minute)
	}

	passwordVerificationToken := passwordVerificationTokenEntity.NewPasswordVerificationTokeFromDatabase(
		id,
		userID,
		createdAt,
		expiresAt,
	)

	return passwordVerificationToken
}

func IngestPasswordVerificationToken(
	t *testing.T,
	testPasswordVerificationToken CreatePasswordVerificationToken,
	savePasswordVerificationToken func(ctx context.Context, passwordVerificationToken passwordVerificationTokenEntity.PasswordVerificationToken) error,
) {
	if (testPasswordVerificationToken != CreatePasswordVerificationToken{}) {
		err := savePasswordVerificationToken(
			context.Background(),
			GeneratePasswordVerificationTokenEntity(t, testPasswordVerificationToken),
		)
		if err != nil {
			t.Fatal(err)
		}
	}
}
