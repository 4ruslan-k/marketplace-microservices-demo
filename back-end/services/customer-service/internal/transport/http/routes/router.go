package routes

import (
	"customer_service/config"
	applicationServices "customer_service/internal/application/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func NewRouter(
	handler *gin.Engine,
	u applicationServices.CustomerApplicationService,
	logger zerolog.Logger,
	config *config.Config,
) {

	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	handler.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })
}
