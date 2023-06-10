package dto

type CreateUserNotificationInput struct {
	UserID             string      `json:"user_id" validate:"required"`
	NotificationTypeID string      `json:"notification_id" validate:"required"`
	Data               interface{} `json:"data"`
}
