package routes

import (
	"authentication/config"
	applicationServices "authentication/internal/services"
	controllers "authentication/internal/transport/http/controllers"
	middlewares "authentication/internal/transport/http/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRouter(
	handler *gin.Engine,
	u applicationServices.UserApplicationService,
	m middlewares.Middlewares,
	logger zerolog.Logger,
	config *config.Config,
	mb *mongo.Database,
) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	handler.Use(m.Session.Apply)
	handler.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	controllers.SetupSocialLogin(mb, config)
	r := controllers.NewUserControllers(u, logger, config)

	v1 := handler.Group("/v1")

	// users
	v1.POST("/users", r.CreateUser)
	v1.GET("/users/:userID", r.GetUserByID)
	v1.PATCH("/users/:userID", r.UpdateUser)
	v1.DELETE("/users/:userID", r.DeleteUser)
	v1.GET("/users/me", r.GetCurrentUser)
	v1.GET("/users/me/internal", r.GetCurrentUserInternal)

	// auth
	v1.POST("/auth/login", r.LoginWithEmailAndPassword)
	v1.POST("/auth/login/mfa/totp", r.LoginWithTotpCode)
	v1.GET("/auth/logout", r.Logout)
	v1.PATCH("/auth/me/change_password", r.ChangeCurrentPassword)
	v1.PUT("/auth/me/mfa/totp", r.GenerateTotpSetup)
	v1.PATCH("/auth/me/mfa/totp/enable", r.EnableTotpMfa)
	v1.PATCH("/auth/me/mfa/totp/disable", r.DisableTotpMfa)

	// auth/social
	v1.GET("/auth/social/:provider/callback", r.SocialLoginCallback)
	v1.GET("/auth/social/:provider", r.SocialLogin)

}
