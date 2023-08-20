package routes

import (
	"net/http"
	"notification_service/config"
	applicationServices "notification_service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	controllers "notification_service/internal/transport/http/controllers"
	socketServer "notification_service/internal/transport/http/socketio"
)

func NewRouter(
	handler *gin.Engine,
	n applicationServices.NotificationApplicationService,
	logger zerolog.Logger,
	config *config.Config,
	socketServer *socketServer.SocketIOServer,
) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	handler.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	handler.GET("/socket.io/*any", gin.WrapH(socketServer.Server))
	handler.POST("/socket.io/*any", gin.WrapH(socketServer.Server))

	r := controllers.NewNotificationController(n, logger, config)

	v1 := handler.Group("/v1")
	v1.GET("/users/me/notifications", r.GetNotificationsByNotificationID)
	v1.PATCH("/users/me/notifications/view", r.ViewAllNotifications)
	v1.DELETE("/users/me/notifications/:notificationId", r.DeleteUserNotification)
	v1.PATCH("/users/me/notifications/:notificationId/view", r.ViewNotification)
}
