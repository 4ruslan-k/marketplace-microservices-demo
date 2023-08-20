package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"cart_service/pkg/httpserver"
)

func run() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	var httpServer *httpserver.Server

	userMessageHandlers, productMessageHandlers, httpServer, err := buildDependencies()
	// TODO: defer pg

	userMessageHandlers.Init()
	productMessageHandlers.Init()

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
