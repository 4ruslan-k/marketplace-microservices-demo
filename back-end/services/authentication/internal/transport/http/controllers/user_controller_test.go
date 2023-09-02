package controllers_test

import (
	"authentication/config"
	applicationServiceMock "authentication/internal/mocks/services"
	applicationService "authentication/internal/services"
	"authentication/internal/transport/http/middlewares"
	routes "authentication/internal/transport/http/routes"
	"errors"
	"io"

	sessionMock "authentication/internal/mocks/sessions"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dto "authentication/internal/services/dto"
)

func doRequest(t *testing.T, req *http.Request) (string, int) {
	t.Helper()

	client := http.Client{}

	res, err := client.Do(req)
	require.NoError(t, err)

	defer func() {
		require.NoError(t, res.Body.Close())
	}()

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	return string(body), res.StatusCode
}

func NewServer(
	t *testing.T, sessionManager *sessionMock.MockSessionManager,
	applicationServiceMock *applicationServiceMock.MockUserApplicationService,
) *httptest.Server {
	t.Helper()

	logger := zerolog.Logger{}
	handler := gin.New()

	sessionStore := cookie.NewStore([]byte("secret"))

	session := middlewares.NewSession(sessionStore)

	m := middlewares.Middlewares{
		Session: session,
	}
	config := &config.Config{}

	routes.NewRouter(handler, applicationServiceMock, m, logger, config, sessionStore, sessionManager)

	server := httptest.NewServer(http.Handler(handler))

	return server
}

func TestUserControllers_GetUserByID(t *testing.T) {
	t.Parallel()
	type fields struct {
		ApplicationService applicationService.UserApplicationService
		Logger             zerolog.Logger
		Config             *config.Config
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionManagerMock := sessionMock.NewMockSessionManager(ctrl)
	applicationServiceMock := applicationServiceMock.NewMockUserApplicationService(ctrl)

	server := NewServer(t, sessionManagerMock, applicationServiceMock)

	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	type args struct {
		userID string
	}

	type want struct {
		statusCode int
		body       string
	}

	testCases := []struct {
		name         string
		args         args
		want         want
		prepareMocks func()
	}{
		{
			name: "success",
			args: args{userID: "abc"},
			want: want{body: `{
				"user": {
					"id": "abc",
					"name": "abc",
					"email": "abc",
					"isMfaEnabled": false
				}
			  }`, statusCode: http.StatusOK},
			prepareMocks: func() {
				userOutput := dto.UserOutput{
					ID:           "abc",
					Name:         "abc",
					Email:        "abc",
					IsMfaEnabled: false,
				}
				sessionManagerMock.EXPECT().GetUserID(gomock.Any()).Return("abc")
				applicationServiceMock.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(&userOutput, nil)
			},
		},
		{
			name: "success",
			args: args{userID: "abz"},
			want: want{body: `{
				"user": {
					"id": "567",
					"name": "568",
					"email": "16",
					"isMfaEnabled": true
				}
			  }`, statusCode: http.StatusOK},
			prepareMocks: func() {
				userOutput := dto.UserOutput{
					ID:           "567",
					Name:         "568",
					Email:        "16",
					IsMfaEnabled: true,
				}
				sessionManagerMock.EXPECT().GetUserID(gomock.Any()).Return("abz")
				applicationServiceMock.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(&userOutput, nil)
			},
		},
		{
			name: "error: unexpected error",
			args: args{userID: "abz"},
			want: want{body: `{
				"message": "Something went wrong",
				"success": false
			  }`, statusCode: http.StatusInternalServerError},
			prepareMocks: func() {
				sessionManagerMock.EXPECT().GetUserID(gomock.Any()).Return("abz")
				applicationServiceMock.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
			},
		},
		{
			name: "error: unauthorized user",
			args: args{userID: "abz"},
			want: want{body: `{
				"message": "You are not authorized to fetch this user",
				"success": false
			  }`, statusCode: http.StatusUnauthorized},
			prepareMocks: func() {
				sessionManagerMock.EXPECT().GetUserID(gomock.Any()).Return("")
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if tc.prepareMocks != nil {
				tc.prepareMocks()
			}

			requestURL := fmt.Sprintf("%s/v1/users/%s", server.URL, tc.args.userID)
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, http.NoBody)

			require.NoError(t, err)

			body, statusCode := doRequest(t, req)
			require.NoError(t, err)

			assert.Equal(t, tc.want.statusCode, statusCode)
			assert.JSONEq(t, tc.want.body, string(body))
		})
	}
}
