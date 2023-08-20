package controllers

import (
	"analytics_service/config"
	applicationServices "analytics_service/internal/application/services"
	userEntity "analytics_service/internal/domain/entities/user"
	"errors"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/gin-gonic/gin"

	domainDto "analytics_service/internal/application/dto"
	httpDto "analytics_service/internal/transport/http/dto"
	httpErrors "analytics_service/pkg/errors/http"
)

type UserControllers struct {
	ApplicationService applicationServices.UserApplicationService
	Logger             zerolog.Logger
	Config             *config.Config
}

type UserOutput struct {
	User *domainDto.UserOutput `json:"user"`
}

func NewUserControllers(
	appService applicationServices.UserApplicationService,
	logger zerolog.Logger,
	config *config.Config,
) *UserControllers {
	return &UserControllers{
		ApplicationService: appService,
		Logger:             logger,
		Config:             config,
	}
}

func (r *UserControllers) GetUserByID(c *gin.Context) {
	userID, found := c.Params.Get("userID")
	if found == false {
		httpErrors.RespondWithError(c, errors.New("user_id parameter not found"))
		return
	}
	user, err := r.ApplicationService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}
	handleResponseWithBody(c, UserOutput{User: user})
}

// Creates a user from user input
func (h *UserControllers) CreateUser(c *gin.Context) {
	var createUserInput httpDto.CreateUserInput
	if err := c.ShouldBindJSON(&createUserInput); err != nil {
		httpErrors.BadRequest(c, err.Error())
		return
	}

	_, err := h.ApplicationService.CreateUser(
		c.Request.Context(),
		userEntity.CreateUserParams{
			Name:  createUserInput.Name,
			Email: createUserInput.Email,
		},
	)
	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}
	handleSuccessResponse(c, http.StatusCreated, "created")
}
