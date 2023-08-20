package controllers

import (
	"customer_service/config"
	applicationServices "customer_service/internal/application/services"
	"errors"

	"github.com/rs/zerolog"

	"github.com/gin-gonic/gin"

	domainDto "customer_service/internal/application/dto"
	httpErrors "shared/errors/http"
)

type CustomerControllers struct {
	ApplicationService applicationServices.CustomerApplicationService
	Logger             zerolog.Logger
	Config             *config.Config
}

type CustomerOutput struct {
	Customer *domainDto.CustomerOutput `json:"customer"`
}

func NewCustomerControllers(
	appService applicationServices.CustomerApplicationService,
	logger zerolog.Logger,
	config *config.Config,
) *CustomerControllers {
	return &CustomerControllers{
		ApplicationService: appService,
		Logger:             logger,
		Config:             config,
	}
}

func (r *CustomerControllers) GetCustomerByID(c *gin.Context) {
	customerID, found := c.Params.Get("customerID")
	if found == false {
		httpErrors.RespondWithError(c, errors.New("customer_id parameter not found"))
		return
	}
	customer, err := r.ApplicationService.GetCustomerByID(c.Request.Context(), customerID)
	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}
	handleResponseWithBody(c, CustomerOutput{Customer: customer})
}
