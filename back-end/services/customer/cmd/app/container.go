package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"customer/config"
	pgStorage "shared/storage/pg"

	httpServ "customer/internal/transport/http"

	applicationServices "customer/internal/services"
	"customer/pkg/httpserver"

	domainServices "customer/internal/domain/services"
	customerRepo "customer/internal/repositories/customer/pg"
	nats "shared/messaging/nats"

	messaging "customer/internal/transport/messaging"
)

func buildDependencies() (messaging.UserMessagingHandlers, *httpserver.Server, error) {

	logger := zerolog.New(os.Stdout)
	config, err := config.NewConfig()
	if err != nil {
		return nil, nil, err
	}

	pg := pgStorage.NewClient(logger, pgStorage.Config{DSN: config.PgSDN})

	customerRepo := customerRepo.NewCustomerRepository(pg, logger)
	customerDomainService := domainServices.NewCustomerService(logger, customerRepo)
	nats := nats.NewNatsClient()

	customerAppService := applicationServices.NewCustomerApplicationService(customerRepo, logger, customerDomainService)

	userMessagingHandlers := messaging.NewCustomerMessagingHandlers(nats, customerAppService, logger)

	httpServer := httpServ.NewHTTPServer(customerAppService, gin.New(), logger, config, pg)

	return userMessagingHandlers, httpServer, nil
}
