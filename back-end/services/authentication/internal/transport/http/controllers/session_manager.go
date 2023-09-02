package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type SessionManager interface {
	GetUserID(*gin.Context) string
}

type sessionManager struct{}

func NewSessionManager() SessionManager {
	return &sessionManager{}
}

func (s *sessionManager) GetUserID(c *gin.Context) string {
	session := sessions.Default(c)
	userID := getUserIDFromSession(session)
	return userID
}
