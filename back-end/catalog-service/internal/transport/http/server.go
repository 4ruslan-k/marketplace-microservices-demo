package httpserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"catalog_service/config"
	"catalog_service/pkg/httpserver"

	applicationServices "catalog_service/internal/application/services"
	routes "catalog_service/internal/transport/http/routes"
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
