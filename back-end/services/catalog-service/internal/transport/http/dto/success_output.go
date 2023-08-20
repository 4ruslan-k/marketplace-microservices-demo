package dto

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuccessOutput struct {
	Message string `json:"message,omitempty" example:"message"`
	Success bool   `json:"success" example:"true"`
}

func HandleSuccessResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, SuccessOutput{Message: msg, Success: true})
}

func HandleOkResponse(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusOK, SuccessOutput{Message: "ok", Success: true})
}

func HandleResponseWithBody(c *gin.Context, body any) {
	c.AbortWithStatusJSON(http.StatusOK, body)
}
