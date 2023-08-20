package httpserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"catalog/config"
	"catalog/pkg/httpserver"

	applicationServices "catalog/internal/services"
	routes "catalog/internal/transport/http/routes"
)

func NewHTTPServer(
	productApplicationService applicationServices.ProductApplicationService,
	handler *gin.Engine,
	logger zerolog.Logger,
	config *config.Config,
) *httpserver.Server {
	routes.NewRouter(handler, productApplicationService, logger, config)
	logger.Info().Msg(fmt.Sprintf("Listening on %s port", config.HTTP.Port))
	return httpserver.New(http.Handler(handler), httpserver.Port(config.HTTP.Port))
}
