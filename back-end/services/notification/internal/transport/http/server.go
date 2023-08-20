package httpserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"

	"notification/config"
	applicationServices "notification/internal/services"
	routes "notification/internal/transport/http/routes"
	socketService "notification/internal/transport/http/socketio"
	"notification/pkg/httpserver"
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
