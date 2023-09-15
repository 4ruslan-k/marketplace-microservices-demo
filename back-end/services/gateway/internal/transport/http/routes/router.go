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

	controllers "gateway/internal/transport/http/controllers"
)

func NewRouter(
	handler *gin.Engine,
	u applicationServices.UserApplicationService,
	m middlewares.Middlewares,
	logger zerolog.Logger,
	config *config.Config,
) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	handler.Use(requestid.New())
	handler.Use(m.GetAuthenticationInfo.Apply)
	handler.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	authenticate := m.RequireAuthentication.Apply

	notificationServiceURL := config.NotificationServiceURL
	notifProxy := controllers.ReverseProxy(notificationServiceURL)

	authServiceURL := config.AuthenticationServiceURL
	authServiceProxy := controllers.ReverseProxy(authServiceURL)

	catalogServiceURL := config.CatalogServiceURL
	catalogServiceProxy := controllers.ReverseProxy(catalogServiceURL)

	cartServiceURL := config.CartServiceURL
	cartServiceProxy := controllers.ReverseProxy(cartServiceURL)

	chatServiceWebsocketURL := config.ChatsServiceWebsocketURL
	chatServiceWebsocketProxy := controllers.ReverseProxy(chatServiceWebsocketURL)

	chatServiceURL := config.ChatsServiceURL
	chatServiceProxy := controllers.ReverseProxy(chatServiceURL)

	// declare before CORS
	// notification service websocket
	handler.GET("/socket.io/*any", authenticate, notifProxy)
	handler.POST("/socket.io/*any", authenticate, notifProxy)

	// chat service websocket
	handler.GET("/chat/socket.io/*any", authenticate, chatServiceWebsocketProxy)
	handler.POST("/chat/socket.io/*any", authenticate, chatServiceWebsocketProxy)

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		config.MarketplaceAppUrl,
		config.AccountsAppURL,
		config.SwaggerEditorDomain,
		config.SwaggerUIDomain,
	}
	corsConfig.AllowCredentials = true

	handler.Use(cors.New(corsConfig))
	v1 := handler.Group("/v1")

	rateLimit := m.RateLimiter.Apply

	// users
	v1.POST("users", authServiceProxy)
	v1.GET("/users/:userID", authenticate, authServiceProxy)
	v1.PATCH("/users/:userID", authenticate, authServiceProxy)
	v1.DELETE("/users/:userID", authenticate, authServiceProxy)
	v1.GET("/users/me", authServiceProxy)

	// chat
	v1.GET("/chat/messages", authenticate, chatServiceProxy)

	// auth
	v1.POST("/auth/login", rateLimit(10), authServiceProxy)
	v1.POST("/auth/login/mfa/totp", rateLimit(10), authServiceProxy)
	v1.GET("/auth/logout", authServiceProxy)
	v1.PATCH("/auth/me/change_password", rateLimit(5), authenticate, authServiceProxy)
	v1.PUT("/auth/me/mfa/totp", rateLimit(5), authenticate, authServiceProxy)
	v1.PATCH("/auth/me/mfa/totp/enable", rateLimit(10), authenticate, authServiceProxy)
	v1.PATCH("/auth/me/mfa/totp/disable", rateLimit(10), authenticate, authServiceProxy)

	// auth/social
	v1.GET("/auth/social/:provider/callback", rateLimit(10), authServiceProxy)
	v1.GET("/auth/social/:provider", rateLimit(10), authServiceProxy)

	// user notifications
	v1.GET("/users/me/notifications", authenticate, notifProxy)
	v1.PATCH("/users/me/notifications/view", authenticate, notifProxy)
	v1.DELETE("/users/me/notifications/:notificationId", authenticate, notifProxy)
	v1.PATCH("/users/me/notifications/:notificationId/view", authenticate, notifProxy)

	// products
	v1.POST("/products", authenticate, catalogServiceProxy)
	v1.GET("/products", authenticate, catalogServiceProxy)
	v1.GET("/products/:productID", authenticate, catalogServiceProxy)
	v1.DELETE("/products/:productID", authenticate, catalogServiceProxy)
	v1.PATCH("/products/:productID", authenticate, catalogServiceProxy)

	// cart
	v1.PATCH("/cart/products", authenticate, cartServiceProxy)
	v1.GET("/cart", authenticate, cartServiceProxy)
}
