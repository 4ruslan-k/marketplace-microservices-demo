package user_test

import (
	user "authentication/internal/domain/entities/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserEntity_NewUser(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name   string
		args   user.CreateUserParams
		expErr error
	}{
		{
			name: "valid_input",
			args: user.CreateUserParams{
				Name:          "name",
				Email:         "example@gmail.coms",
				Password:      "example@gmail.coms",
				SocialAccount: nil,
			},
		},
		{
			name: "bad_email_format",
			args: user.CreateUserParams{
				Name:          "name",
				Email:         "example.com",
				Password:      "password",
				SocialAccount: nil,
			},
			expErr: user.ErrInvalidEmailFormat,
		},
		{
			name: "empty_password",
			args: user.CreateUserParams{
				Name:          "name",
				Email:         "example@gmail.coms",
				Password:      "",
				SocialAccount: nil,
			},
			expErr: user.ErrInvalidPassword,
		},
	}

	for _, tCase := range testCases {
		tCase := tCase
		t.Run(tCase.name, func(t *testing.T) {
			t.Parallel()
			user, err := user.NewUser(tCase.args)
			if tCase.expErr != nil {
				require.Error(t, err)
				require.EqualError(t, tCase.expErr, err.Error())
			} else {
				require.NoError(t, err)
				assert.NotNil(t, user)
			}
		})
	}
}
