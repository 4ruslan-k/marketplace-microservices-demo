package soccialaccount_test

import (
	socialAccountEntity "authentication/internal/domain/entities/social_account"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSocialAccountEntity_NewSocialAccount(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name   string
		args   socialAccountEntity.CreateSocialAccountParams
		expErr error
	}{
		{
			name: "valid_input",
			args: socialAccountEntity.CreateSocialAccountParams{
				ID:       "id",
				Name:     "name",
				Email:    "example@gmail.coms",
				Provider: "provider",
			},
		},
		{
			name: "bad_email_format",
			args: socialAccountEntity.CreateSocialAccountParams{
				ID:       "id",
				Name:     "name",
				Email:    "example.com",
				Provider: "google",
			},
			expErr: socialAccountEntity.ErrInvalidEmailFormat,
		},
		{
			name: "empty_provider",
			args: socialAccountEntity.CreateSocialAccountParams{
				ID:       "id",
				Name:     "name",
				Email:    "example@gmail.com",
				Provider: "",
			},
			expErr: socialAccountEntity.ErrInvalidProvider,
		},
	}

	for _, tCase := range testCases {
		tCase := tCase
		t.Run(tCase.name, func(t *testing.T) {
			t.Parallel()
			socialAccount, err := socialAccountEntity.NewSocialAccount(tCase.args)
			if tCase.expErr != nil {
				require.Error(t, err)
				require.EqualError(t, tCase.expErr, err.Error())
			} else {
				require.NoError(t, err)
				assert.NotNil(t, socialAccount)
			}
		})
	}
}
