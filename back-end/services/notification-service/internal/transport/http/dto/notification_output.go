package dto

import (
	domainDto "notification_service/internal/services/dto"
)

type NotificationOutput struct {
	Notifications []domainDto.NotificationOutput `json:"notifications"`
}
