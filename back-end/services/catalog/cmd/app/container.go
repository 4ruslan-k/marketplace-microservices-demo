package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"catalog/config"
	nats "shared/messaging/nats"
	pgStorage "shared/storage/pg"

	applicationServices "catalog/internal/services"

	repository "catalog/internal/repositories/product/pg"
	httpServ "catalog/internal/transport/http"
	"catalog/pkg/httpserver"
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
