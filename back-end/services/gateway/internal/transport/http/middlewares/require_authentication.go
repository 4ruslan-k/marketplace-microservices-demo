package middlewares

import (
	applicationServices "gateway/internal/domain/application-services"
	httpErrors "shared/errors/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type RequireAuthentication struct {
	logger zerolog.Logger
}

func (r RequireAuthentication) Apply(c *gin.Context) {
	userFromContext, exists := c.Get("user")
	if exists {
		user := userFromContext.(applicationServices.User)
		if user.ID == "" {
			httpErrors.BadRequest(c, "please login")
			return
		}
	} else {
		httpErrors.BadRequest(c, "please login")
		return
	}
	c.Next()
}

func NewRequireAuthentication(
	logger zerolog.Logger,
) *RequireAuthentication {
	return &RequireAuthentication{logger}
}
