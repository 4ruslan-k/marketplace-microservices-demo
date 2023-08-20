package dto

import (
	domainDto "analytics_service/internal/application/dto"
)

type UserOutput struct {
	User *domainDto.UserOutput `json:"user"`
}
