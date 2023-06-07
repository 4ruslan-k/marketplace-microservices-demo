package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"catalog_service/config"
)

func Run() {
	_, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("config.NewConfig")
	}

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
