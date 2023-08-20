package app

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"notification_service/config"
	"notification_service/pkg/httpserver"
	pgStorage "notification_service/pkg/storage/pg"

	httpServ "notification_service/internal/transport/http"
	socketServer "notification_service/internal/transport/http/socketio"

	applicationServices "notification_service/internal/application/services"
	domainServices "notification_service/internal/domain/services"
	notificationInfraRepository "notification_service/internal/infrastructure/repositories/pg/notification"
	userInfraRepository "notification_service/internal/infrastructure/repositories/pg/user"
	nats "notification_service/pkg/messaging/nats"

	messaging "notification_service/internal/transport/messaging"
)

func buildDependencies() (
	messaging.UserMessagingHandlers,
	messaging.NotificationMessagingHandlers,
	*socketServer.SocketIOServer,
	*httpserver.Server,
	error,
) {
	logger := zerolog.New(os.Stdout)
	config, err := config.NewConfig()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	pg := pgStorage.NewClient(logger, config)
	nats := nats.NewNatsClient()
	userRepo := userInfraRepository.NewUserRepository(pg, logger)
	notificationRepo := notificationInfraRepository.NewNotificationRepository(pg, logger)

	userDomainService := domainServices.NewUserService(logger, userRepo)

	userAppService := applicationServices.NewUserApplicationService(userRepo, logger, userDomainService)
	notificationAppService := applicationServices.NewNotificationApplicationService(notificationRepo, logger, nats)

	userMessageHandlers := messaging.NewUserMessagingHandlers(nats, userAppService, logger)
	socketServer := socketServer.NewSocketIOServer(logger)
	notificationMessageHandlers := messaging.NewNotificationMessagingHandlers(nats, notificationAppService, logger, socketServer)
	httpServer := httpServ.NewHTTPServer(notificationAppService, gin.New(), logger, config, pg, socketServer)
	return userMessageHandlers, notificationMessageHandlers, socketServer, httpServer, nil
}
