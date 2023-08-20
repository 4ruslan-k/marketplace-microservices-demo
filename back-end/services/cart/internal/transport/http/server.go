package httpserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"

	"cart/config"
	controllers "cart/internal/transport/http/controllers"
	routes "cart/internal/transport/http/routes"

	"cart/pkg/httpserver"
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
