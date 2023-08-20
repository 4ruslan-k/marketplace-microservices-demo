package dto

import (
	domainDto "authentication_service/internal/application/dto"
)

type LoginOutput struct {
	LoginOutput domainDto.LoginOutput `json:"user"`
}
