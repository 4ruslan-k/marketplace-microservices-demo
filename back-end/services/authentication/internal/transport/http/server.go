package httpserver

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"authentication/config"
	"authentication/pkg/httpserver"

	applicationServices "authentication/internal/services"
	controller "authentication/internal/transport/http/controllers"
	middlewares "authentication/internal/transport/http/middlewares"
	routes "authentication/internal/transport/http/routes"
)

func NewHTTPServer(
	userApplicationService applicationServices.UserApplicationService,
	handler *gin.Engine,
	m middlewares.Middlewares,
	logger zerolog.Logger,
	config *config.Config,
	sessionsStore sessions.Store,
) *httpserver.Server {
	sessionManager := controller.NewSessionManager()
	routes.NewRouter(handler, userApplicationService, m, logger, config, sessionsStore, sessionManager)
	logger.Info().Msg(fmt.Sprintf("Listening on %s port", config.HTTP.Port))
	return httpserver.New(http.Handler(handler), httpserver.Port(config.HTTP.Port))
}
