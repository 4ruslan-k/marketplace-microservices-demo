package middlewares

import (
	"authentication_service/config"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/mongo/mongodriver"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
)

type Session struct {
	Apply gin.HandlerFunc
}

func NewSession(sessionStore sessions.Store) *Session {

	sessionMiddleware := sessions.Sessions("session", sessionStore)
	return &Session{sessionMiddleware}
}

func NewSessionStore(m *mongo.Database, config *config.Config) sessions.Store {
	sessionsCollection := m.Collection("sessions")
	store := mongodriver.NewStore(sessionsCollection, 3600, true, []byte(config.SessionSecret))
	return store
}
