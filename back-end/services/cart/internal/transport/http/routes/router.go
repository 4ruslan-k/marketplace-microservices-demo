package routes

import (
	"cart/config"
	controllers "cart/internal/transport/http/controllers"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func NewRouter(
	handler *gin.Engine,
	p *controllers.ProductController,
	logger zerolog.Logger,
	config *config.Config,
) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	handler.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	v1 := handler.Group("/v1")

	v1.PATCH("/cart/products", p.UpdateProductsInCart)
	v1.GET("/cart", p.GetCart)

}
