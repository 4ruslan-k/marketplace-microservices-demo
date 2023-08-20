package applicationservices

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	domainDto "notification_service/internal/application/dto"
	notificationEntity "notification_service/internal/domain/entities/notification"
	"notification_service/internal/ports/repositories"
	customErrors "notification_service/pkg/errors"
	nats "shared/messaging/nats"

	"github.com/rs/zerolog"
)

var (
	ErrInvalidUserID        = customErrors.NewIncorrectInputError("invalid_input", "User ID not provided")
	ErrNotificationNotFound = customErrors.NewIncorrectInputError("not_found", "User notification not found")
)

var _ NotificationApplicationService = (*notificationApplicationService)(nil)

type notificationApplicationService struct {
	notificationRepository repositories.NotificationsRepository
	logger                 zerolog.Logger
	natsClient             nats.NatsClient
}

type UserNotificationRefetchEvent struct {
	UserID string `json:"id"`
}

type NotificationApplicationService interface {
	CreateUserNotification(
		ctx context.Context,
		createNotificationParams notificationEntity.CreateUserNotificationParams,
	) error
	GetNotificationsByUserID(ctx context.Context, userID string) (domainDto.NotificationListOutput, error)
	ViewNotification(ctx context.Context, userID string, userNotificationID string) error
	ViewAllNotifications(ctx context.Context, userID string) error
	DeleteUserNotification(ctx context.Context, userID string, userNotificationID string) error
}

func NewNotificationApplicationService(
	notificationRepository repositories.NotificationsRepository,
	logger zerolog.Logger,
	natsClient nats.NatsClient,
) notificationApplicationService {
	return notificationApplicationService{notificationRepository, logger, natsClient}
}

const (
	refetchUserNotificationSubject = "notifications.refetch"
)

func (n notificationApplicationService) CreateUserNotification(
	ctx context.Context,
	createUserNotificationParams notificationEntity.CreateUserNotificationParams,
) error {

	notification := notificationEntity.NotificationByTypeIds[createUserNotificationParams.NotificationTypeID]
	createUserNotificationParams.Title = notification.TitleTemplate()
	createUserNotificationParams.Message = notification.MessageTemplate()

	userNotification, err := notificationEntity.NewUserNotification(createUserNotificationParams)
	if err != nil {
		return err
	}
	err = n.notificationRepository.CreateUserNotification(ctx, userNotification)
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(UserNotificationRefetchEvent{
		UserID: createUserNotificationParams.UserID,
	})
	if err != nil {
		return err
	}

	n.natsClient.PublishMessageEphemeral("notifications.refetch", string(bytes))

	return nil
}

func (n notificationApplicationService) GetNotificationsByUserID(
	ctx context.Context,
	userID string,
) (domainDto.NotificationListOutput, error) {
	if userID == "" {
		return domainDto.NotificationListOutput{}, ErrInvalidUserID
	}
	notifications, err := n.notificationRepository.GetByUserID(ctx, userID)
	notificationOutput := make([]domainDto.NotificationOutput, 0, len(notifications))
	if err != nil {
		return domainDto.NotificationListOutput{}, err
	}
	for _, notification := range notifications {
		var viewedAt *time.Time
		if !notification.ViewedAt().IsZero() {
			viewedAtValue := notification.ViewedAt()
			viewedAt = &viewedAtValue
		}
		notificationOutput = append(notificationOutput, domainDto.NotificationOutput{
			ID:                 notification.ID(),
			NotificationTypeID: notification.NotificationTypeID(),
			ViewedAt:           viewedAt,
			CreatedAt:          notification.CreatedAt(),
			Message:            notification.Message(),
			Title:              notification.Title(),
		})
	}
	return domainDto.NotificationListOutput{Notifications: notificationOutput}, nil
}

func (n notificationApplicationService) ViewNotification(ctx context.Context, userID string, userNotificationID string) error {

	userNotification, err := n.notificationRepository.GetByUserIDAndUserNotificationID(ctx, userID, userNotificationID)
	if err != nil {
		return fmt.Errorf("notificationApplicationService -> ViewNotification -> GetByUserIDAndUserNotificationID: %w", err)
	}

	if userNotification.IsZero() {
		return ErrNotificationNotFound
	}
	err = n.notificationRepository.MarkUserNotificationViewed(ctx, userID, userNotification.ID())

	if err != nil {
		return fmt.Errorf("notificationApplicationService -> ViewNotification -> MarkUserNotificationViewed: %w", err)
	}

	bytes, err := json.Marshal(UserNotificationRefetchEvent{
		UserID: userID,
	})

	if err != nil {
		return fmt.Errorf("notificationApplicationService -> ViewNotification -> Marshal: %w", err)
	}

	n.natsClient.PublishMessageEphemeral(refetchUserNotificationSubject, string(bytes))

	return nil
}

func (n notificationApplicationService) ViewAllNotifications(ctx context.Context, userID string) error {

	err := n.notificationRepository.MarkAllUserNotificationViewed(ctx, userID)

	if err != nil {
		return fmt.Errorf("notificationApplicationService ViewAllNotifications -> MarkAllUserNotificationViewed: %w", err)
	}

	bytes, err := json.Marshal(UserNotificationRefetchEvent{
		UserID: userID,
	})

	if err != nil {
		return fmt.Errorf("notificationApplicationService ViewAllNotifications -> json.Marshal: %w", err)
	}

	n.natsClient.PublishMessageEphemeral(refetchUserNotificationSubject, string(bytes))

	return nil
}

func (n notificationApplicationService) DeleteUserNotification(ctx context.Context, userID string, userNotificationID string) error {

	err := n.notificationRepository.DeleteUserNotification(ctx, userID, userNotificationID)

	if err != nil {
		return fmt.Errorf("notificationApplicationService DeleteUserNotification -> NewDelete: %w", err)
	}

	bytes, err := json.Marshal(UserNotificationRefetchEvent{
		UserID: userID,
	})
	if err != nil {
		return fmt.Errorf("notificationApplicationService DeleteUserNotification -> json.Marshal: %w", err)
	}

	n.natsClient.PublishMessageEphemeral(refetchUserNotificationSubject, string(bytes))

	return nil
}
