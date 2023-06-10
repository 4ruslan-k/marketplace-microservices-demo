package controllers

import (
	"encoding/json"
	applicationServices "gateway/internal/domain/application-services"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/rs/zerolog"

	"github.com/gin-gonic/gin"
)

type UserHandlers struct {
	ApplicationService applicationServices.UserApplicationService
	Logger             zerolog.Logger
}

type dataHandlers struct {
	logger zerolog.Logger
}

type AuthenticationInfo struct {
	UserID    string `json:"user_id"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	SessionID string `json:"session_id"`
}

func ReverseProxy(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		target, err := url.Parse(target)
		userFromContext, exists := c.Get("user")
		var authenticationInfo AuthenticationInfo
		if exists {
			user := userFromContext.(applicationServices.User)
			userID := user.ID
			ip := c.ClientIP()
			userAgent := c.Request.UserAgent()
			authenticationInfo = AuthenticationInfo{
				UserID:    userID,
				IP:        ip,
				UserAgent: userAgent,
				SessionID: user.SessionID,
			}
		}
		if err != nil {
			panic(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(target)
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			authInfoJSON, err := json.Marshal(authenticationInfo)
			if err != nil {
				panic(err)
			}
			req.Header.Add("X-Authentication-Info", string(authInfoJSON))
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
