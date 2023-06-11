package routes

import (
	"cart_service/config"
	applicationServices "cart_service/internal/application/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func NewRouter(
	handler *gin.Engine,
	n applicationServices.ProductApplicationService,
	logger zerolog.Logger,
	config *config.Config,
) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	handler.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })
}
