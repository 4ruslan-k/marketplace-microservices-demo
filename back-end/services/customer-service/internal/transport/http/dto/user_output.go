package dto

import (
	domainDto "customer_service/internal/application/dto"
)

type CustomerOutput struct {
	User *domainDto.CustomerOutput `json:"customer"`
}
