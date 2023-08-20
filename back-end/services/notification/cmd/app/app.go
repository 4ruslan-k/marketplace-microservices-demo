package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"notification/pkg/httpserver"

	socketServer "notification/internal/transport/http/socketio"
)

func run() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	var httpServer *httpserver.Server
	var socketServ *socketServer.SocketIOServer

	userMessageHandlers, notificationMessageHandlers, socketServer, httpServer, err := buildDependencies()
	// TODO: defer pg

	userMessageHandlers.Init()
	notificationMessageHandlers.Init()
	socketServ = socketServer

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

	defer socketServ.Server.Close()

}
