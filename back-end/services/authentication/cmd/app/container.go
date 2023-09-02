package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"authentication/config"
	"authentication/pkg/httpserver"
	storage "authentication/pkg/storage/mongo"

	domainServices "authentication/internal/domain/services"
	authRepository "authentication/internal/repositories/authentication/mongo"
	userRepository "authentication/internal/repositories/user/mongo"
	applicationServices "authentication/internal/services"
	httpServ "authentication/internal/transport/http"
	middlewares "authentication/internal/transport/http/middlewares"
	nats "shared/messaging/nats"
)

func buildDependencies() (*httpserver.Server, error) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	logger := zerolog.New(os.Stdout)

	config, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	mongo := storage.NewMongoClient(logger, config)

	authenticationRepo := authRepository.NewAuthenticationRepository(mongo, logger)
	userRepo := userRepository.NewUserRepository(mongo, logger)

	authenticationDomainService := domainServices.NewAuthenticationService(logger, authenticationRepo)
	userDomainService := domainServices.NewUserService(logger, authenticationDomainService, userRepo)

	nats := nats.NewNatsClient()

	userApplicationService := applicationServices.NewUserApplicationService(
		userRepo, authenticationRepo, logger, userDomainService, authenticationDomainService, nats)

	sessionStore := middlewares.NewSessionStore(mongo, config)
	session := middlewares.NewSession(sessionStore)

	middlewaresContainer := middlewares.Middlewares{
		Session: session,
	}
	server := httpServ.NewHTTPServer(userApplicationService, gin.New(), middlewaresContainer, logger, config, sessionStore)

	return server, nil
}
