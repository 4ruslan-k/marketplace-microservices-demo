package controllers

import (
	"github.com/gin-gonic/gin"
)

type successResponse struct {
	Message string `json:"message" example:"message"`
}

func handleSuccessResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, successResponse{msg})
}
