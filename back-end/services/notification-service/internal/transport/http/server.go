package httpserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"

	"notification_service/config"
	applicationServices "notification_service/internal/services"
	routes "notification_service/internal/transport/http/routes"
	socketService "notification_service/internal/transport/http/socketio"
	"notification_service/pkg/httpserver"
)

func NewHTTPServer(
	notificationApplicationService applicationServices.NotificationApplicationService,
	handler *gin.Engine,
	logger zerolog.Logger,
	config *config.Config,
	db *bun.DB,
	socketServer *socketService.SocketIOServer,

) *httpserver.Server {
	routes.NewRouter(handler, notificationApplicationService, logger, config, socketServer)
	logger.Info().Msg(fmt.Sprintf("Listening on %s port", config.HTTP.Port))
	return httpserver.New(http.Handler(handler), httpserver.Port(config.HTTP.Port))
}
