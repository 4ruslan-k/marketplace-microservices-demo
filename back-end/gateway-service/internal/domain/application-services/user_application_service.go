package applicationServices

import (
	"encoding/json"
	"gateway/config"

	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type userApplicationService struct {
	logger zerolog.Logger
	config *config.Config
}

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	SessionID string `json:"session_id"`
}

type UserResponse struct {
	User      User   `json:"user"`
	SessionID string `json:"session_id"`
}

type UserApplicationService interface {
	GetCurrentUser(http.Header) (User, string, error)
}

func NewUserApplicationService(
	logger zerolog.Logger,
	config *config.Config,
) userApplicationService {
	return userApplicationService{logger, config}
}

func (u userApplicationService) GetCurrentUser(headers http.Header) (User, string, error) {
	req, err := http.NewRequest("GET", u.config.UserServiceURL+"/v1/users/me/internal", nil)
	if err != nil {
		return User{}, "", err
	}

	req.Header = headers
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return User{}, "", err

	}

	var userResponse UserResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&userResponse)

	if err != nil {
		return User{}, "", errors.Wrap(err, "GetCurrentUser -> decoding user error")
	}

	return userResponse.User, userResponse.SessionID, nil
}
