package httperrors

import (
	"customer_service/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func InternalError(c *gin.Context, message string) {
	httpRespondWithError(c, message, http.StatusInternalServerError)
}

func Unauthorized(c *gin.Context, message string) {
	httpRespondWithError(c, message, http.StatusUnauthorized)
}

func BadRequest(c *gin.Context, message string) {
	httpRespondWithError(c, message, http.StatusBadRequest)
}

func NotFoundRequest(c *gin.Context, message string) {
	httpRespondWithError(c, message, http.StatusNotFound)
}

func RespondWithError(c *gin.Context, err error) {
	log.Debug().Err(err).Msg("RespondWithError")
	errorStruct, ok := err.(errors.CustomError)
	if !ok {
		InternalError(c, "Something went wrong")
		return
	}
	switch errorStruct.ErrorType() {
	case errors.ErrorTypeAuthorization:
		Unauthorized(c, errorStruct.Message())
	case errors.ErrorTypeIncorrectInput:
		BadRequest(c, errorStruct.Message())
	case errors.ErrorTypeNotFound:
		NotFoundRequest(c, errorStruct.Message())
	default:
		InternalError(c, errorStruct.Message())
	}
}

func httpRespondWithError(
	c *gin.Context,
	message string,
	status int,
) {
	resp := ErrorResponse{Message: message, Success: false, httpStatus: status}
	c.JSON(resp.httpStatus, resp)
}

type ErrorResponse struct {
	Message    string `json:"message"`
	Success    bool   `json:"success"`
	httpStatus int
}

func (e ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(e.httpStatus)
	return nil
}
