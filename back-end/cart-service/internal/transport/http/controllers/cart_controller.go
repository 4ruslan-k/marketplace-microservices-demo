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

type UpdateProductsInCartInput struct {
	ProductID string `json:"productId" binding:"required"`
	Quantity  int    `json:"quantity" binding:"gte=0"`
}

func (p *ProductController) UpdateProductsInCart(c *gin.Context) {
	var UpdateProductsInCartInput UpdateProductsInCartInput
	if err := c.ShouldBindJSON(&UpdateProductsInCartInput); err != nil {
		httpErrors.BadRequest(c, err.Error())
		return
	}

	var authInfo AuthInfo
	authValue := c.Request.Header.Get("X-Authentication-Info")
	json.Unmarshal([]byte(authValue), &authInfo)

	err := p.ApplicationService.UpdateProductsInCart(
		c.Request.Context(),
		UpdateProductsInCartInput.ProductID,
		UpdateProductsInCartInput.Quantity,
		authInfo.UserID,
	)
	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}
	handleOkResponse(c)
}

type GetCartInput struct {
	CustomerID string `json:"productId" binding:"required"`
}

func (p *ProductController) GetCart(c *gin.Context) {

	var authInfo AuthInfo
	authValue := c.Request.Header.Get("X-Authentication-Info")
	json.Unmarshal([]byte(authValue), &authInfo)

	cart, err := p.ApplicationService.GetCart(
		c.Request.Context(),
		authInfo.UserID,
	)
	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}
	handleResponseWithBody(c, cart)
}
