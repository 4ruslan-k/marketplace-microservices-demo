package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"authentication_service/config"
	"authentication_service/pkg/httpserver"
	storage "authentication_service/pkg/storage/mongo"

	domainServices "authentication_service/internal/domain/services"
	authRepository "authentication_service/internal/repositories/authentication/mongo"
	userRepository "authentication_service/internal/repositories/user/mongo"
	applicationServices "authentication_service/internal/services"
	httpServ "authentication_service/internal/transport/http"
	middlewares "authentication_service/internal/transport/http/middlewares"
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

	server := httpServ.NewHTTPServer(userApplicationService, gin.New(), middlewaresContainer, logger, config, mongo)

	return server, nil
}
