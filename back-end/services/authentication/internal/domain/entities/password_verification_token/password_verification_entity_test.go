package passwordveificationtoken_test

import (
	passwordVerificationTokenEntity "authentication/internal/domain/entities/password_verification_token"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUserPasswordVerificationTokenEntity_NewUser(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		args        passwordVerificationTokenEntity.CreatePasswordVerificationToken
		currentTime time.Time
		isExpired   bool
	}{
		{
			name: "not expired (valid)",
			args: passwordVerificationTokenEntity.CreatePasswordVerificationToken{
				ID:                 "someIdToken",
				UserID:             "userIdTest",
				CurrentTime:        time.Now(),
				ExpirationDuration: time.Minute * 10,
			},
			currentTime: time.Now(),
			isExpired:   false,
		},
		{
			name: "expired (not valid)",
			args: passwordVerificationTokenEntity.CreatePasswordVerificationToken{
				ID:                 "someIdTokenOne",
				UserID:             "userIdTest",
				CurrentTime:        time.Now(),
				ExpirationDuration: time.Minute * 10,
			},
			currentTime: time.Now().Add(time.Minute * 10),
			isExpired:   true,
		},
	}

	for _, tCase := range testCases {
		tCase := tCase
		t.Run(tCase.name, func(t *testing.T) {
			t.Parallel()
			passwordVerificationToken := passwordVerificationTokenEntity.NewPasswordVerificationToken(tCase.args)
			require.Equal(t, tCase.args.ID, passwordVerificationToken.ID())
			require.Equal(t, tCase.args.UserID, passwordVerificationToken.UserID())
			require.Equal(t, tCase.args.CurrentTime, passwordVerificationToken.CreatedAt())
			require.Equal(t, tCase.args.CurrentTime.Add(tCase.args.ExpirationDuration), passwordVerificationToken.ExpiresAt())
			isExpired := passwordVerificationToken.HasExpired(tCase.currentTime)
			require.Equal(t, tCase.isExpired, isExpired)
		})
	}
}
