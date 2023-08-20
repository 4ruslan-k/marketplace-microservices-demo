package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"analytics/config"
	pgStorage "shared/storage/pg"

	httpServ "analytics/internal/transport/http"

	applicationServices "analytics/internal/services"
	"analytics/pkg/httpserver"

	domainServices "analytics/internal/domain/services"
	repository "analytics/internal/repositories/user/pg"
	nats "shared/messaging/nats"

	messaging "analytics/internal/transport/messaging"
)

func buildDependencies() (messaging.UserMessagingHandlers, *httpserver.Server, error) {

	logger := zerolog.New(os.Stdout)
	config, err := config.NewConfig()
	if err != nil {
		return nil, nil, err
	}

	pg := pgStorage.NewClient(logger, pgStorage.Config{DSN: config.PgSDN})

	userRepo := repository.NewUserRepository(pg, logger)
	userDomainService := domainServices.NewUserService(logger, userRepo)
	nats := nats.NewNatsClient()

	userAppService := applicationServices.NewUserApplicationService(userRepo, logger, userDomainService)

	userMessagingHandlers := messaging.NewUserMessagingHandlers(nats, userAppService, logger)

	httpServer := httpServ.NewHTTPServer(userAppService, gin.New(), logger, config, pg)

	return userMessagingHandlers, httpServer, nil
}
