package routes

import (
	"catalog_service/config"
	applicationServices "catalog_service/internal/application/services"
	controllers "catalog_service/internal/transport/http/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func NewRouter(
	handler *gin.Engine,
	u applicationServices.ProductApplicationService,
	logger zerolog.Logger,
	config *config.Config,
) {
	handler.Use(otelgin.Middleware(config.App.Name))
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	r := controllers.NewProductController(u, logger, config)

	v1 := handler.Group("/v1")

	// products
	v1.POST("/products", r.CreateProduct)
	v1.GET("/products", r.GetProducts)
	v1.GET("/products/:productID", r.GetProductByID)
	v1.DELETE("/products/:productID", r.DeleteProductByID)
}
