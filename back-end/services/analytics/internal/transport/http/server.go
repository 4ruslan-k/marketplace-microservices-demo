package httpserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"

	"analytics/config"
	"analytics/pkg/httpserver"

	applicationServices "analytics/internal/services"
	routes "analytics/internal/transport/http/routes"
)

func NewHTTPServer(
	userApplicationService applicationServices.UserApplicationService,
	handler *gin.Engine,
	logger zerolog.Logger,
	config *config.Config,
	db *bun.DB,
) *httpserver.Server {
	routes.NewRouter(handler, userApplicationService, logger, config)
	logger.Info().Msg(fmt.Sprintf("Listening on %s port", config.HTTP.Port))
	return httpserver.New(http.Handler(handler), httpserver.Port(config.HTTP.Port))
}
