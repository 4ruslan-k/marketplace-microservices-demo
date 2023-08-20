package dto

import (
	domainDto "notification/internal/services/dto"
)

type NotificationOutput struct {
	Notifications []domainDto.NotificationOutput `json:"notifications"`
}
