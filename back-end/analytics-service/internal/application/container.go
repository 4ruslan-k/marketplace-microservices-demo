package app

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"analytics_service/config"
	pgStorage "analytics_service/pkg/storage/pg"

	httpServ "analytics_service/internal/transport/http"

	applicationServices "analytics_service/internal/application/services"
	"analytics_service/pkg/httpserver"

	domainServices "analytics_service/internal/domain/services"
	userInfraRepository "analytics_service/internal/infrastructure/repositories/pg/user"
	nats "analytics_service/pkg/messaging/nats"

	messaging "analytics_service/internal/transport/messaging"
)

func buildDependencies() (messaging.UserMessagingHandlers, *httpserver.Server, error) {

	logger := zerolog.New(os.Stdout)
	config, err := config.NewConfig()
	if err != nil {
		return nil, nil, err
	}

	pg := pgStorage.NewClient(logger, config)

	userRepo := userInfraRepository.NewUserRepository(pg, logger)
	userDomainService := domainServices.NewUserService(logger, userRepo)
	nats := nats.NewNatsClient()

	userAppService := applicationServices.NewUserApplicationService(userRepo, logger, userDomainService)

	userMessagingHandlers := messaging.NewUserMessagingHandlers(nats, userAppService, logger)

	httpServer := httpServ.NewHTTPServer(userAppService, gin.New(), logger, config, pg)

	return userMessagingHandlers, httpServer, nil
}
