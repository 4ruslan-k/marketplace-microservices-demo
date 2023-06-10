package dto

import "time"

type NotificationOutput struct {
	ID                 string     `json:"id"`
	Title              string     `json:"title"`
	Message            string     `json:"message"`
	NotificationTypeID string     `json:"notificationTypeId"`
	ViewedAt           *time.Time `json:"viewedAt,omitempty"`
	CreatedAt          time.Time  `json:"createdAt,omitempty"`
}
