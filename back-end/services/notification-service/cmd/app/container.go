package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"notification_service/config"
	"notification_service/pkg/httpserver"
	pgStorage "shared/storage/pg"

	httpServ "notification_service/internal/transport/http"
	socketServer "notification_service/internal/transport/http/socketio"

	domainServices "notification_service/internal/domain/services"
	notificationRepository "notification_service/internal/repositories/notification/pg"
	userRepository "notification_service/internal/repositories/user/pg"
	applicationServices "notification_service/internal/services"
	nats "shared/messaging/nats"

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
	pg := pgStorage.NewClient(logger, pgStorage.Config{DSN: config.PgSDN})
	nats := nats.NewNatsClient()
	userRepo := userRepository.NewUserRepository(pg, logger)
	notificationRepo := notificationRepository.NewNotificationRepository(pg, logger)

	userDomainService := domainServices.NewUserService(logger, userRepo)

	userAppService := applicationServices.NewUserApplicationService(userRepo, logger, userDomainService)
	notificationAppService := applicationServices.NewNotificationApplicationService(notificationRepo, logger, nats)

	userMessageHandlers := messaging.NewUserMessagingHandlers(nats, userAppService, logger)
	socketServer := socketServer.NewSocketIOServer(logger)
	notificationMessageHandlers := messaging.NewNotificationMessagingHandlers(nats, notificationAppService, logger, socketServer)
	httpServer := httpServ.NewHTTPServer(notificationAppService, gin.New(), logger, config, pg, socketServer)
	return userMessageHandlers, notificationMessageHandlers, socketServer, httpServer, nil
}
