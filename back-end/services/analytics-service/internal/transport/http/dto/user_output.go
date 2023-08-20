package dto

import (
	domainDto "analytics_service/internal/services/dto"
)

type UserOutput struct {
	User *domainDto.UserOutput `json:"user"`
}
