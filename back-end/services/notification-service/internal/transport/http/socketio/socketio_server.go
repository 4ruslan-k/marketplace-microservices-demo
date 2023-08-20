package socketserver

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"github.com/rs/zerolog"
)

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

type Connections struct {
	sync.RWMutex
	connections map[string][]socketio.Conn
}

type SocketIOServer struct {
	Server      *socketio.Server
	Connections Connections
	logger      zerolog.Logger
}

type AuthInfo struct {
	UserID string `json:"user_id"`
}

func (s *SocketIOServer) SendEvent(userId string) {
	s.Connections.RLock()
	defer s.Connections.RUnlock()
	soList, ok := s.Connections.connections[userId]
	for _, so := range soList {
		if ok {
			s.logger.Info().Msgf("sending message to user: %s", userId)

			so.Emit("notifications_update", "Hello from server!")
		} else {
			s.logger.Info().Msgf("NOT sending message to user: %s not found", userId)

		}
	}

}

func (s *SocketIOServer) initialize() {
	s.Server.OnConnect("/", func(so socketio.Conn) error {
		var authInfo AuthInfo
		authValue := so.RemoteHeader().Get("X-Authentication-Info")
		json.Unmarshal([]byte(authValue), &authInfo)
		s.logger.Info().Msgf("New client connected: %s", so.ID())
		s.Connections.Lock()
		defer s.Connections.Unlock()
		so.SetContext(authInfo)
		s.Connections.connections[authInfo.UserID] = append(s.Connections.connections[authInfo.UserID], so)

		return nil
	})

	s.Server.OnEvent("/", "message", func(so socketio.Conn, msg string) {
		s.logger.Info().Msgf("Received message: %s", msg)
		so.Emit("reply", "Received message: "+time.Now().Format(time.RFC3339))
	})

	s.Server.OnEvent("/", "reply", func(so socketio.Conn, msg string) {
		s.logger.Info().Msgf("Received reply: %s", msg)
		so.Emit("reply", "Received reply: "+time.Now().Format(time.RFC3339))
	})

	s.Server.OnDisconnect("/", func(so socketio.Conn, reason string) {
		s.logger.Info().Msgf("Client disconnected: %s", so.ID())
		s.logger.Info().Msgf("Disconnect reason: %s", reason)

	})

}

func NewSocketIOServer(logger zerolog.Logger) *SocketIOServer {
	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})

	socketServer := &SocketIOServer{Server: server, logger: logger, Connections: Connections{connections: make(map[string][]socketio.Conn)}}
	socketServer.initialize()

	go func() {
		if err := server.Serve(); err != nil {
			logger.Error().Err(err).Msg("socketio listen error")
		}
	}()

	return socketServer
}
