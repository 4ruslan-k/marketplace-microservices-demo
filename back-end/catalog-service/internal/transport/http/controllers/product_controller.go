package controllers

import (
	"catalog_service/config"
	applicationServices "catalog_service/internal/application/services"
	dto "catalog_service/internal/transport/http/dto"

	"github.com/rs/zerolog"

	"github.com/gin-gonic/gin"

	httpErrors "catalog_service/pkg/errors/http"

	productEntity "catalog_service/internal/domain/entities/product"
)

type ProductController struct {
	ApplicationService applicationServices.ProductApplicationService
	Logger             zerolog.Logger
	Config             *config.Config
}

func NewProductController(
	appService applicationServices.ProductApplicationService,
	logger zerolog.Logger,
	config *config.Config,
) *ProductController {
	return &ProductController{
		ApplicationService: appService,
		Logger:             logger,
		Config:             config,
	}
}

// Creates product
func (h *ProductController) CreateProduct(c *gin.Context) {
	var createProductInput dto.CreateProductInput
	if err := c.ShouldBindJSON(&createProductInput); err != nil {
		httpErrors.BadRequest(c, err.Error())
		return
	}

	err := h.ApplicationService.CreateProduct(
		c.Request.Context(),
		productEntity.CreateProductParams{Name: createProductInput.Name},
	)

	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}
	dto.HandleOkResponse(c)
}
