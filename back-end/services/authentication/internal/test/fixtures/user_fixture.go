package fixtures

import (
	mfaSettingsEntity "authentication/internal/domain/entities/mfa_settings"
	socialAccountEntity "authentication/internal/domain/entities/social_account"
	userEntity "authentication/internal/domain/entities/user"
	domainService "authentication/internal/domain/services"
	"context"
	"testing"
	"time"

	"github.com/jaswdr/faker"
)

var fake = faker.New()

func GenerateUUID() string {
	UUID := fake.UUID().V4()
	return UUID
}

func GenerateRandomEmail() string {
	email := fake.Internet().Email()
	return email
}

func GenerateRandomName() string {
	email := fake.Person().Name()
	return email
}

func GenerateRandomPassword() string {
	password := fake.Internet().Password()
	return password
}

type CreateTestUser struct {
	ID           string
	Email        string
	Password     string
	Name         string
	PasswordHash string
	TotpSecret   string
	IsMfaEnabled bool
}

func GenerateUserEntity(t *testing.T, c CreateTestUser) userEntity.User {
	t.Helper()
	id := c.ID

	if id == "" {
		id = fake.UUID().V4()
	}

	email := c.Email
	if email == "" {
		email = GenerateRandomEmail()
	}

	name := c.Name
	if name == "" {
		name = GenerateRandomName()
	}

	passwordHash := c.PasswordHash
	var err error
	if passwordHash == "" {
		password := c.Password
		if password == "" {
			password = GenerateRandomPassword()
		}
		passwordHash, err = domainService.GetPasswordHashValue(password)
		if err != nil {
			t.Fatal(err)
		}
	}
	mfaSettings := mfaSettingsEntity.NewMfaSettingsFromDatabase(
		c.IsMfaEnabled,
		c.TotpSecret,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		t.Fatal(err)
	}

	user, err := userEntity.NewUserFromDatabase(
		id,
		email,
		name,
		passwordHash,
		time.Now(),
		nil,
		[]socialAccountEntity.SocialAccount{},
		mfaSettings,
	)
	if err != nil {
		t.Fatal(err)
	}
	return *user
}

func IngestUser(
	t *testing.T,
	testUser CreateTestUser,
	createUser func(ctx context.Context, user userEntity.User) (string, error),
) {
	if (testUser != CreateTestUser{}) {
		_, err := createUser(
			context.Background(),
			GenerateUserEntity(t, testUser),
		)
		if err != nil {
			t.Fatal(err)
		}
	}
}
