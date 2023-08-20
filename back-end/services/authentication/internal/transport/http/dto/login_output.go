package dto

import (
	domainDto "authentication/internal/services/dto"
)

type LoginOutput struct {
	LoginOutput domainDto.LoginOutput `json:"user"`
}
