package app

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"catalog_service/config"
	pgStorage "catalog_service/pkg/storage/pg"

	applicationServices "catalog_service/internal/application/services"

	pgRepositories "catalog_service/internal/infrastructure/repositories/pg"
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
	pgConn := pgStorage.NewClient(logger, conf.PgDSN)
	productRepo := pgRepositories.NewProductRepository(pgConn, logger)
	productAppService := applicationServices.NewProductApplicationService(productRepo, logger)
	server := httpServ.NewHTTPServer(productAppService, gin.New(), logger, conf)

	return server, nil
}
