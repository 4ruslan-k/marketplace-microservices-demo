package repositories

import (
	"context"
	notificationEntity "notification_service/internal/domain/entities/notification"
)

type NotificationsRepository interface {
	CreateUserNotification(ctx context.Context, userNotification notificationEntity.UserNotification) error
	GetByUserID(ctx context.Context, userID string) ([]notificationEntity.UserNotification, error)
	GetByUserIDAndUserNotificationID(ctx context.Context, userID string, userNotificationID string) (notificationEntity.UserNotification, error)
	MarkUserNotificationViewed(ctx context.Context, userID string, userNotificationID string) error
	MarkAllUserNotificationViewed(ctx context.Context, userID string) error
	DeleteUserNotification(ctx context.Context, userID string, userNotificationID string) error
}
