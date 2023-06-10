package fixtures

import (
	"context"
	notificationEntity "notification_service/internal/domain/entities/notification"
	"testing"
	"time"
)

type CreateTestUserNotification struct {
	ID                 string
	UserID             string
	NotificationTypeID string
	TypeID             string
	Data               interface{}
	CreatedAt          time.Time
	ViewedAt           time.Time
	UpdatedAt          time.Time
	Message            string
	Title              string
}

func GenerateUserNotificationEntity(t *testing.T, c CreateTestUserNotification) notificationEntity.UserNotification {
	t.Helper()
	id := c.ID

	if id == "" {
		id = GenerateUUID()
	}

	userID := c.UserID
	if userID == "" {
		userID = GenerateUUID()
	}

	notificationTypeID := c.NotificationTypeID
	if notificationTypeID == "" {
		notificationTypeID = GenerateUUID()
	}

	createdAt := c.CreatedAt

	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	viewedAt := c.ViewedAt

	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	var err error

	if err != nil {
		t.Fatal(err)
	}

	notification := notificationEntity.NewUserNotificationFromDatabase(
		id,
		userID,
		notificationTypeID,
		c.Data,
		createdAt,
		viewedAt,
		c.UpdatedAt,
		c.Message,
		c.Title,
	)
	if err != nil {
		t.Fatal(err)
	}
	return notification
}

func InsertUserNotification(
	t *testing.T,
	testNotification CreateTestUserNotification,
	createUserNotification func(ctx context.Context, testNotification notificationEntity.UserNotification) error,
) {
	if (testNotification != CreateTestUserNotification{}) {
		err := createUserNotification(
			context.Background(),
			GenerateUserNotificationEntity(t, testNotification),
		)
		if err != nil {
			t.Fatal(err)
		}
	}
}
