package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"catalog_service/config"
	nats "shared/messaging/nats"
	pgStorage "shared/storage/pg"

	applicationServices "catalog_service/internal/services"

	repository "catalog_service/internal/repositories/product/pg"
	httpServ "catalog_service/internal/transport/http"
	"catalog_service/pkg/httpserver"
)

func buildDependencies() (*httpserver.Server, error) {
	conf, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout)
	pgConn := pgStorage.NewClient(logger, pgStorage.Config{DSN: conf.PgDSN})
	natsClient := nats.NewNatsClient()
	productRepo := repository.NewProductRepository(pgConn, logger)
	productAppService := applicationServices.NewProductApplicationService(productRepo, logger, natsClient)
	server := httpServ.NewHTTPServer(productAppService, gin.New(), logger, conf)

	return server, nil
}
