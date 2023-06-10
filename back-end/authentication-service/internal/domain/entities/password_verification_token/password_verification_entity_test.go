package passwordveificationtoken_test

import (
	passwordVerificationTokenEntity "authentication_service/internal/domain/entities/password_verification_token"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUserPasswordVerificationTokenEntity_NewUser(t *testing.T) {

	cases := []struct {
		name        string
		in          passwordVerificationTokenEntity.CreatePasswordVerificationToken
		currentTime time.Time
		isExpired   bool
	}{
		{
			name: "not expired (valid)",
			in: passwordVerificationTokenEntity.CreatePasswordVerificationToken{
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
			in: passwordVerificationTokenEntity.CreatePasswordVerificationToken{
				ID:                 "someIdTokenOne",
				UserID:             "userIdTest",
				CurrentTime:        time.Now(),
				ExpirationDuration: time.Minute * 10,
			},
			currentTime: time.Now().Add(time.Minute * 10),
			isExpired:   true,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			passwordVerificationToken := passwordVerificationTokenEntity.NewPasswordVerificationToken(tCase.in)
			require.Equal(t, tCase.in.ID, passwordVerificationToken.ID())
			require.Equal(t, tCase.in.UserID, passwordVerificationToken.UserID())
			require.Equal(t, tCase.in.CurrentTime, passwordVerificationToken.CreatedAt())
			require.Equal(t, tCase.in.CurrentTime.Add(tCase.in.ExpirationDuration), passwordVerificationToken.ExpiresAt())
			isExpired := passwordVerificationToken.HasExpired(tCase.currentTime)
			require.Equal(t, tCase.isExpired, isExpired)
		})
	}
}
