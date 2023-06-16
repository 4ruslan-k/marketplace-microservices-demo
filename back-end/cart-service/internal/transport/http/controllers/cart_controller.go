package controllers

import (
	"cart_service/config"
	applicationServices "cart_service/internal/application/services"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	httpErrors "cart_service/pkg/errors/http"
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

type AuthInfo struct {
	UserID string `json:"user_id"`
}

type AddProductToCartInput struct {
	ProductID string `json:"productId" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required"`
}

func (p *ProductController) AddProductToCart(c *gin.Context) {
	var addProductToCartInput AddProductToCartInput
	if err := c.ShouldBindJSON(&addProductToCartInput); err != nil {
		httpErrors.BadRequest(c, err.Error())
		return
	}

	var authInfo AuthInfo
	authValue := c.Request.Header.Get("X-Authentication-Info")
	json.Unmarshal([]byte(authValue), &authInfo)

	err := p.ApplicationService.AddProductToCart(
		c.Request.Context(),
		addProductToCartInput.ProductID,
		addProductToCartInput.Quantity,
		authInfo.UserID,
	)
	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}
	handleOkResponse(c)
}