package dto

import (
	domainDto "notification_service/internal/application/dto"
)

type NotificationOutput struct {
	Notifications []domainDto.NotificationOutput `json:"notifications"`
}
