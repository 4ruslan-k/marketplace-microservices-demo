package httpserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"

	"authentication_service/config"
	"authentication_service/pkg/httpserver"

	applicationServices "authentication_service/internal/application/services"
	middlewares "authentication_service/internal/transport/http/middlewares"
	routes "authentication_service/internal/transport/http/routes"
)

func NewHTTPServer(
	userApplicationService applicationServices.UserApplicationService,
	handler *gin.Engine,
	m middlewares.Middlewares,
	logger zerolog.Logger,
	config *config.Config,
	mb *mongo.Database,
) *httpserver.Server {
	routes.NewRouter(handler, userApplicationService, m, logger, config, mb)
	logger.Info().Msg(fmt.Sprintf("Listening on %s port", config.HTTP.Port))
	return httpserver.New(http.Handler(handler), httpserver.Port(config.HTTP.Port))
}
