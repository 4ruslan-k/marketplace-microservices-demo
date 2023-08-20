package notification

import (
	"time"

	customErrors "shared/errors"

	"github.com/google/uuid"
)

var (
	ErrUserIDEmpty              = customErrors.NewIncorrectInputError("invalid_input", "User ID is not set")
	ErrNotificationTypeIDEmpty  = customErrors.NewIncorrectInputError("invalid_input", "Notification type ID is not set")
	ErrNotificationMessageEmpty = customErrors.NewIncorrectInputError("invalid_input", "Notification message is not set")
	ErrNotificationTitleEmpty   = customErrors.NewIncorrectInputError("invalid_input", "Notification title is not set")
)

type UserNotification struct {
	id                 string
	userID             string
	notificationTypeID string
	title              string
	message            string
	createdAt          time.Time
	viewedAt           time.Time
	updatedAt          time.Time
	data               interface{}
}

type CreateUserNotificationParams struct {
	UserID             string
	NotificationTypeID string
	Data               interface{}
	Message            string
	Title              string
}

func NewUserNotification(createNotificationParams CreateUserNotificationParams) (UserNotification, error) {
	if createNotificationParams.UserID == "" {
		return UserNotification{}, ErrUserIDEmpty
	}

	if createNotificationParams.NotificationTypeID == "" {
		return UserNotification{}, ErrNotificationTypeIDEmpty
	}

	if createNotificationParams.Message == "" {
		return UserNotification{}, ErrNotificationMessageEmpty
	}

	if createNotificationParams.Title == "" {
		return UserNotification{}, ErrNotificationTitleEmpty
	}

	notification := UserNotification{
		id:                 uuid.NewString(),
		notificationTypeID: createNotificationParams.NotificationTypeID,
		userID:             createNotificationParams.UserID,
		data:               createNotificationParams.Data,
		createdAt:          time.Now(),
		message:            createNotificationParams.Message,
		title:              createNotificationParams.Title,
	}
	return notification, nil
}

func NewUserNotificationFromDatabase(
	id string,
	userID string,
	notificationTypeID string,
	data interface{},
	createdAt time.Time,
	viewedAt time.Time,
	updatedAt time.Time,
	message string,
	title string,
) UserNotification {
	userNotification := UserNotification{id: id,
		userID:             userID,
		notificationTypeID: notificationTypeID,
		data:               data,
		createdAt:          createdAt,
		viewedAt:           viewedAt,
		updatedAt:          updatedAt,
		message:            message,
		title:              title,
	}
	return userNotification
}

func (d UserNotification) ID() string {
	return d.id
}

func (d UserNotification) UserID() string {
	return d.userID
}

func (d UserNotification) CreatedAt() time.Time {
	return d.createdAt
}

func (d UserNotification) ViewedAt() time.Time {
	return d.viewedAt
}

func (d UserNotification) UpdatedAt() time.Time {
	return d.updatedAt
}

func (d UserNotification) Data() interface{} {
	return d.data
}

func (d UserNotification) NotificationTypeID() string {
	return d.notificationTypeID
}

func (d UserNotification) Message() string {
	return d.message
}

func (d UserNotification) Title() string {
	return d.title
}

func (d UserNotification) IsZero() bool {
	return d == UserNotification{}
}
