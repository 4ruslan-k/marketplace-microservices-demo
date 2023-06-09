package app

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"customer_service/config"
	pgStorage "customer_service/pkg/storage/pg"

	httpServ "customer_service/internal/transport/http"

	applicationServices "customer_service/internal/application/services"
	"customer_service/pkg/httpserver"

	domainServices "customer_service/internal/domain/services"
	customerInfraRepository "customer_service/internal/infrastructure/repositories/pg/customer"
	nats "customer_service/pkg/messaging/nats"

	messaging "customer_service/internal/transport/messaging"
)

func buildDependencies() (messaging.UserMessagingHandlers, *httpserver.Server, error) {

	logger := zerolog.New(os.Stdout)
	config, err := config.NewConfig()
	if err != nil {
		return nil, nil, err
	}

	pg := pgStorage.NewClient(logger, config)

	customerRepo := customerInfraRepository.NewCustomerRepository(pg, logger)
	customerDomainService := domainServices.NewCustomerService(logger, customerRepo)
	nats := nats.NewNatsClient()

	customerAppService := applicationServices.NewCustomerApplicationService(customerRepo, logger, customerDomainService)

	userMessagingHandlers := messaging.NewCustomerMessagingHandlers(nats, customerAppService, logger)

	httpServer := httpServ.NewHTTPServer(customerAppService, gin.New(), logger, config, pg)

	return userMessagingHandlers, httpServer, nil
}
