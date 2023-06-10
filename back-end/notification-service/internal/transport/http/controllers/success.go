package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuccessOutput struct {
	Message string `json:"message,omitempty" example:"message"`
	Success bool   `json:"success" example:"true"`
}

func handleSuccessResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, SuccessOutput{Message: msg, Success: true})
}

func handleOkResponse(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusOK, SuccessOutput{Message: "ok", Success: true})
}

func handleResponseWithBody(c *gin.Context, body any) {
	c.AbortWithStatusJSON(http.StatusOK, body)
}
