package app

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gateway/config"
	"gateway/pkg/httpserver"

	applicationServices "gateway/internal/domain/application-services"
	middlewares "gateway/internal/transport/http/middlewares"
	routes "gateway/internal/transport/http/routes"

	redisPool "gateway/pkg/storage/redis"
)

func NewHTTPServer(
	userApplicationService applicationServices.UserApplicationService,
	handler *gin.Engine,
	m middlewares.Middlewares,
	logger zerolog.Logger,
	config *config.Config,
) *httpserver.Server {
	routes.NewRouter(handler, userApplicationService, m, logger, config)
	port := config.HTTP.Port
	logger.Info().Msg(fmt.Sprintf("Listening on %s port", port))
	return httpserver.New(http.Handler(handler), httpserver.Port(port))
}

func buildDependencies() (*httpserver.Server, error) {

	logger := zerolog.New(os.Stdout)
	config, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	redis := redisPool.NewRedisPool(config.RedisAddress)

	rateLimiter := middlewares.NewRateLimiter(redis)
	requireAuthentication := middlewares.NewRequireAuthentication(logger)
	userApplicationService := applicationServices.NewUserApplicationService(logger, config)

	getAuthenticationInfo := middlewares.NewGetAuthenticationInfo(logger, userApplicationService)

	middlewaresContainer := middlewares.Middlewares{
		RequireAuthentication: requireAuthentication,
		GetAuthenticationInfo: getAuthenticationInfo,
		RateLimiter:           rateLimiter,
	}

	httpServer := NewHTTPServer(userApplicationService, gin.New(), middlewaresContainer, logger, config)

	return httpServer, err
}

// Run creates objects via constructors.
func Run() {
	_, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("config.NewConfig")
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	httpServer, err := buildDependencies()
	if err != nil {
		log.Panic().Err(err).Msg("c.Invoke")
	}
	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info().Msg("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		log.Error().Err(err).Msg("app - Run - httpServer.Notify")
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		log.Error().Err(err).Msg("app - Run - httpServer.Shutdown")
	}

}
