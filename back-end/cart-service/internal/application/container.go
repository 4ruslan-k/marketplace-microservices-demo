package app

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"cart_service/config"
	"cart_service/pkg/httpserver"
	pgStorage "cart_service/pkg/storage/pg"

	httpServ "cart_service/internal/transport/http"

	applicationServices "cart_service/internal/application/services"
	domainServices "cart_service/internal/domain/services"
	userInfraRepository "cart_service/internal/infrastructure/repositories/pg/customer"
	productInfraRepository "cart_service/internal/infrastructure/repositories/pg/product"
	nats "cart_service/pkg/messaging/nats"

	messaging "cart_service/internal/transport/messaging"
)

func buildDependencies() (
	messaging.UserMessagingHandlers,
	messaging.ProductMessagingHandlers,
	*httpserver.Server,
	error,
) {
	logger := zerolog.New(os.Stdout)
	config, err := config.NewConfig()
	if err != nil {
		return nil, nil, nil, err
	}
	pg := pgStorage.NewClient(logger, config)
	nats := nats.NewNatsClient()
	userRepo := userInfraRepository.NewCustomerRepository(pg, logger)
	productRepo := productInfraRepository.NewProductRepository(pg, logger)

	userDomainService := domainServices.NewCustomerService(logger, userRepo)

	userApplicationService := applicationServices.NewCustomerApplicationService(userRepo, logger, userDomainService)
	productAppService := applicationServices.NewProductApplicationService(productRepo, logger, nats)

	userMessageHandlers := messaging.NewCustomerMessagingHandlers(nats, userApplicationService, logger)
	productMessageHandlers := messaging.NewProductMessagingHandlers(nats, productAppService, logger)
	httpServer := httpServ.NewHTTPServer(productAppService, gin.New(), logger, config, pg)
	return userMessageHandlers, productMessageHandlers, httpServer, nil
}
