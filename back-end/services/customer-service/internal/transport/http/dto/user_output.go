package dto

import (
	domainDto "customer_service/internal/services/dto"
)

type CustomerOutput struct {
	User *domainDto.CustomerOutput `json:"customer"`
}
