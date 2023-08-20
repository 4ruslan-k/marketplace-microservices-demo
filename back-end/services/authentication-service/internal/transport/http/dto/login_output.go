package dto

import (
	domainDto "authentication_service/internal/services/dto"
)

type LoginOutput struct {
	LoginOutput domainDto.LoginOutput `json:"user"`
}
