package dto

import (
	domainDto "analytics/internal/services/dto"
)

type UserOutput struct {
	User *domainDto.UserOutput `json:"user"`
}
