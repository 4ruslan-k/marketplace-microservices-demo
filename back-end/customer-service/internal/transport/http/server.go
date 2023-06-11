package httpserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"

	"customer_service/config"
	"customer_service/pkg/httpserver"

	applicationServices "customer_service/internal/application/services"
	routes "customer_service/internal/transport/http/routes"
)

func NewHTTPServer(
	CustomerApplicationService applicationServices.CustomerApplicationService,
	handler *gin.Engine,
	logger zerolog.Logger,
	config *config.Config,
	db *bun.DB,
) *httpserver.Server {
	routes.NewRouter(handler, CustomerApplicationService, logger, config)
	logger.Info().Msg(fmt.Sprintf("Listening on %s port", config.HTTP.Port))
	return httpserver.New(http.Handler(handler), httpserver.Port(config.HTTP.Port))
}
