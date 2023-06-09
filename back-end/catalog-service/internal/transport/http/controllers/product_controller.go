package controllers

import (
	"catalog_service/config"
	applicationServices "catalog_service/internal/application/services"
	dto "catalog_service/internal/transport/http/dto"
	"errors"

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

// Creates a product
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

// Fetches a list of products
func (h *ProductController) GetProducts(c *gin.Context) {
	products, err := h.ApplicationService.GetProducts(c.Request.Context())

	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}
	dto.HandleResponseWithBody(c, dto.NewProductListOutputFromEntities(products))
}

// Fetches a product
func (h *ProductController) GetProductByID(c *gin.Context) {
	productID, found := c.Params.Get("productID")
	if found == false {
		httpErrors.RespondWithError(c, errors.New("product id parameter not found"))
		return
	}
	product, err := h.ApplicationService.GetProductByID(c.Request.Context(), productID)

	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}
	dto.HandleResponseWithBody(c, dto.NewProductOutputFromEntity(product))
}

// Deletes a product
func (h *ProductController) DeleteProductByID(c *gin.Context) {
	productID, found := c.Params.Get("productID")
	if found == false {
		httpErrors.RespondWithError(c, errors.New("product id parameter not found"))
		return
	}
	err := h.ApplicationService.DeleteProductByID(c.Request.Context(), productID)

	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}
	dto.HandleOkResponse(c)
}
