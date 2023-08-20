package dto

import (
	domainDto "authentication/internal/services/dto"
)

type UserOutput struct {
	User *domainDto.UserOutput `json:"user"`
}

type UserWithSessionIDOutput struct {
	User      *domainDto.UserOutput `json:"user"`
	SessionID string                `json:"session_id"`
}
