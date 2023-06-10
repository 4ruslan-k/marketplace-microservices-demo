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

// CreateProduct creates a product
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

// GetProducts fetches a list of products
func (h *ProductController) GetProducts(c *gin.Context) {
	products, err := h.ApplicationService.GetProducts(c.Request.Context())

	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}
	dto.HandleResponseWithBody(c, dto.NewProductListOutputFromEntities(products))
}

// GetProductByID Fetches a product"
func (h *ProductController) GetProductByID(c *gin.Context) {
	var params struct {
		ProductID string `uri:"productID" binding:"required,uuid"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		httpErrors.BadRequest(c, err.Error())
		return
	}
	product, err := h.ApplicationService.GetProductByID(c.Request.Context(), params.ProductID)

	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}
	dto.HandleResponseWithBody(c, dto.NewProductOutputFromEntity(product))
}

// DeleteProductByID deletes a product
func (h *ProductController) DeleteProductByID(c *gin.Context) {
	var params struct {
		ProductID string `uri:"productID" binding:"required,uuid"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		httpErrors.BadRequest(c, err.Error())
		return
	}
	err := h.ApplicationService.DeleteProductByID(c.Request.Context(), params.ProductID)

	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}
	dto.HandleOkResponse(c)
}
