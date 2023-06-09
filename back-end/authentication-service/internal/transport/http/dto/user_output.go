package dto

import (
	domainDto "authentication_service/internal/application/dto"
)

type UserOutput struct {
	User *domainDto.UserOutput `json:"user"`
}

type UserWithSessionIDOutput struct {
	User      *domainDto.UserOutput `json:"user"`
	SessionID string                `json:"session_id"`
}
