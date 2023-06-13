package httpserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"

	"cart_service/config"
	controllers "cart_service/internal/transport/http/controllers"
	routes "cart_service/internal/transport/http/routes"

	"cart_service/pkg/httpserver"
)

func NewHTTPServer(
	productController *controllers.ProductController,
	handler *gin.Engine,
	logger zerolog.Logger,
	config *config.Config,
	db *bun.DB,

) *httpserver.Server {
	routes.NewRouter(handler, productController, logger, config)
	logger.Info().Msg(fmt.Sprintf("Listening on %s port", config.HTTP.Port))
	return httpserver.New(http.Handler(handler), httpserver.Port(config.HTTP.Port))
}
