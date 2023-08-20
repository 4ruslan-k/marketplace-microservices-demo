package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"cart/config"
	"cart/pkg/httpserver"
	pgStorage "shared/storage/pg"

	httpServ "cart/internal/transport/http"

	cartInfraRepository "cart/internal/repositories/cart/pg"
	userInfraRepository "cart/internal/repositories/customer/pg"
	productInfraRepository "cart/internal/repositories/product/pg"
	applicationServices "cart/internal/services"
	nats "shared/messaging/nats"

	messaging "cart/internal/transport/messaging"

	controllers "cart/internal/transport/http/controllers"
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
	pg := pgStorage.NewClient(logger, pgStorage.Config{DSN: config.PgSDN})
	nats := nats.NewNatsClient()
	userRepo := userInfraRepository.NewCustomerRepository(pg, logger)
	productRepo := productInfraRepository.NewProductRepository(pg, logger)
	cartRepo := cartInfraRepository.NewCartRepository(pg, logger)

	userApplicationService := applicationServices.NewCustomerApplicationService(userRepo, logger)
	productAppService := applicationServices.NewProductApplicationService(productRepo, cartRepo, logger, nats)

	productController := controllers.NewProductController(productAppService, logger, config)

	userMessageHandlers := messaging.NewCustomerMessagingHandlers(nats, userApplicationService, logger)
	productMessageHandlers := messaging.NewProductMessagingHandlers(nats, productAppService, logger)
	httpServer := httpServ.NewHTTPServer(productController, gin.New(), logger, config, pg)
	return userMessageHandlers, productMessageHandlers, httpServer, nil
}
