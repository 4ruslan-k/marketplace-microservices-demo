package httperrors

import (
	customErrors "catalog_service/pkg/errors"
	"errors"
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

func TooManyRequests(c *gin.Context) {
	httpRespondWithError(c, "Too Many requests", http.StatusTooManyRequests)
}

func RespondWithError(c *gin.Context, err error) {
	log.Debug().Err(err).Msg("RespondWithError")
	var customErrorStruct customErrors.CustomError
	ok := errors.As(err, &customErrorStruct)
	if !ok {
		InternalError(c, "Something went wrong")
		return
	}
	switch customErrorStruct.ErrorType() {
	case customErrors.ErrorTypeAuthorization:
		Unauthorized(c, customErrorStruct.Message())
	case customErrors.ErrorTypeIncorrectInput:
		BadRequest(c, customErrorStruct.Message())
	case customErrors.ErrorTypeNotFound:
		NotFoundRequest(c, customErrorStruct.Message())
	default:
		InternalError(c, customErrorStruct.Message())
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
