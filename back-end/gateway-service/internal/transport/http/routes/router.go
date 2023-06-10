// Package v1 implements routing paths. Each services in own file.
package routes

import (
	"gateway/config"
	applicationServices "gateway/internal/domain/application-services"
	middlewares "gateway/internal/transport/http/middlewares"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	controllers "gateway/internal/transport/http/controllers"
)

func NewRouter(
	handler *gin.Engine,
	u applicationServices.UserApplicationService,
	m middlewares.Middlewares,
	logger zerolog.Logger,
	config *config.Config,
) {
	handler.Use(otelgin.Middleware(config.App.Name))
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	handler.Use(requestid.New())
	// handler.Use(m.GetAuthenticationInfo.Apply)
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	corsConfig := cors.DefaultConfig()

	authenticate := m.RequireAuthentication.Apply

	corsConfig.AllowOrigins = []string{config.FontendURL, config.SwaggerEditorDomain, config.SwaggerUIDomain}
	corsConfig.AllowCredentials = true

	handler.Use(cors.New(corsConfig))
	v1 := handler.Group("/v1")

	rateLimit := m.RateLimiter.Apply

	userServiceURL := config.UserServiceURL
	userProxy := controllers.ReverseProxy(userServiceURL)

	catalogServiceURL := config.CatalogServiceURL
	catalogServiceProxy := controllers.ReverseProxy(catalogServiceURL)

	// users
	v1.POST("", userProxy)
	v1.GET("/users/:userID", authenticate, userProxy)
	v1.PATCH("/users/:userID", authenticate, userProxy)
	v1.DELETE("/users/:userID", authenticate, userProxy)
	v1.GET("/users/me", userProxy)

	// auth
	v1.POST("/auth/login", rateLimit(10), userProxy)
	v1.POST("/auth/login/mfa/totp", rateLimit(10), userProxy)
	v1.GET("/auth/logout", userProxy)
	v1.PATCH("/auth/me/change_password", rateLimit(5), authenticate, userProxy)
	v1.PUT("/auth/me/mfa/totp", rateLimit(5), authenticate, userProxy)
	v1.PATCH("/auth/me/mfa/totp/enable", rateLimit(10), authenticate, userProxy)
	v1.PATCH("/auth/me/mfa/totp/disable", rateLimit(10), authenticate, userProxy)

	// auth/social
	v1.GET("/auth/social/:provider/callback", rateLimit(10), userProxy)
	v1.GET("/auth/social/:provider", rateLimit(10), userProxy)

	// products
	v1.POST("/products", catalogServiceProxy)
	v1.GET("/products", catalogServiceProxy)
	v1.GET("/products/:productID", catalogServiceProxy)
	v1.DELETE("/products/:productID", catalogServiceProxy)
}
