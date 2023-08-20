package user_test

import (
	user "authentication/internal/domain/entities/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserEntity_NewUser(t *testing.T) {

	cases := []struct {
		name   string
		in     user.CreateUserParams
		expErr error
	}{
		{
			name: "valid_input",
			in: user.CreateUserParams{
				Name:          "name",
				Email:         "example@gmail.coms",
				Password:      "example@gmail.coms",
				SocialAccount: nil,
			},
		},
		{
			name: "bad_email_format",
			in: user.CreateUserParams{
				Name:          "name",
				Email:         "example.com",
				Password:      "password",
				SocialAccount: nil,
			},
			expErr: user.ErrInvalidEmailFormat,
		},
		{
			name: "empty_password",
			in: user.CreateUserParams{
				Name:          "name",
				Email:         "example@gmail.coms",
				Password:      "",
				SocialAccount: nil,
			},
			expErr: user.ErrInvalidPassword,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			user, err := user.NewUser(tCase.in)
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
