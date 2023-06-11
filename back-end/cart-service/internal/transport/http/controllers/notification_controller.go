package controllers

import (
	"cart_service/config"
	applicationServices "cart_service/internal/application/services"

	"github.com/rs/zerolog"
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
