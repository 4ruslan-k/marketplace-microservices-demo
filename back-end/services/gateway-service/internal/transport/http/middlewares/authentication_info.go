package middlewares

import (
	applicationServices "gateway/internal/domain/application-services"
	httpErrors "shared/errors/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type GetAuthenticationInfo struct {
	logger                 zerolog.Logger
	userApplicationService applicationServices.UserApplicationService
}

func (r GetAuthenticationInfo) Apply(c *gin.Context) {
	headers := c.Request.Header
	user, sessionID, err := r.userApplicationService.GetCurrentUser(headers)
	r.logger.Info().Interface("user", user).Str("sessionID", sessionID).Msg("GetAuthenticationInfo -> user")
	if err != nil {
		r.logger.Error().Err(err).Msg("r.userApplicationService.GetCurrentUser -> user")
		httpErrors.BadRequest(c, err.Error())
		return
	}

	user.SessionID = sessionID
	c.Set("user", user)

	c.Next()
}

func NewGetAuthenticationInfo(
	logger zerolog.Logger,
	userApplicationService applicationServices.UserApplicationService,
) *GetAuthenticationInfo {
	return &GetAuthenticationInfo{logger, userApplicationService}
}
