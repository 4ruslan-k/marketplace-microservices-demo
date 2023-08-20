package dto

import (
	domainDto "customer/internal/services/dto"
)

type CustomerOutput struct {
	User *domainDto.CustomerOutput `json:"customer"`
}
