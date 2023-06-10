package fixtures

import (
	"context"
	userEntity "notification_service/internal/domain/entities/user"
	"testing"

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

	var err error

	if err != nil {
		t.Fatal(err)
	}

	user, err := userEntity.NewUserFromDatabase(
		id,
		email,
		name,
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
