package routes

import (
	"catalog/config"
	applicationServices "catalog/internal/services"
	controllers "catalog/internal/transport/http/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func NewRouter(
	handler *gin.Engine,
	u applicationServices.ProductApplicationService,
	logger zerolog.Logger,
	config *config.Config,
) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	handler.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	r := controllers.NewProductController(u, logger, config)

	v1 := handler.Group("/v1")

	// products
	v1.POST("/products", r.CreateProduct)
	v1.GET("/products", r.GetProducts)
	v1.GET("/products/:productID", r.GetProductByID)
	v1.DELETE("/products/:productID", r.DeleteProductByID)
	v1.PATCH("/products/:productID", r.UpdateProductByID)
}
