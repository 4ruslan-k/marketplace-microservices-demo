package httpserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"

	"cart_service/config"
	applicationServices "cart_service/internal/application/services"
	routes "cart_service/internal/transport/http/routes"
	"cart_service/pkg/httpserver"
)

func NewHTTPServer(
	productApplicationService applicationServices.ProductApplicationService,
	handler *gin.Engine,
	logger zerolog.Logger,
	config *config.Config,
	db *bun.DB,

) *httpserver.Server {
	routes.NewRouter(handler, productApplicationService, logger, config)
	logger.Info().Msg(fmt.Sprintf("Listening on %s port", config.HTTP.Port))
	return httpserver.New(http.Handler(handler), httpserver.Port(config.HTTP.Port))
}
