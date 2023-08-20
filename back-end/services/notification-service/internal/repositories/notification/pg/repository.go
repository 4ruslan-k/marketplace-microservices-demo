package repository

import (
	"context"
	"database/sql"
	"fmt"
	notificationEntity "notification_service/internal/domain/entities/notification"
	repositories "notification_service/internal/repositories/notification"
	"time"

	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

type UserNotificationModel struct {
	bun.BaseModel `bun:"table:notifications"`

	ID                 string      `bun:"id,pk"`
	NotificationTypeID string      `bun:"notification_type_id"`
	UserID             string      `bun:"user_id"`
	ViewedAt           time.Time   `bun:"viewed_at,nullzero"`
	Data               interface{} `bun:"data"`
	Message            string      `bun:"message"`
	Title              string      `bun:"title"`
	CreatedAt          time.Time   `bun:"created_at"`
	UpdatedAt          time.Time   `bun:"updated_at,nullzero"`
}

var _ repositories.NotificationsRepository = (*notificationPGRepository)(nil)

type notificationPGRepository struct {
	db     *bun.DB
	logger zerolog.Logger
}

func (n *UserNotificationModel) toDB(un notificationEntity.UserNotification) UserNotificationModel {
	return UserNotificationModel{
		ID:                 un.ID(),
		UserID:             un.UserID(),
		NotificationTypeID: un.NotificationTypeID(),
		Message:            un.Message(),
		Title:              un.Title(),
		Data:               un.Data(),
		CreatedAt:          un.CreatedAt(),
		ViewedAt:           un.ViewedAt(),
	}
}

func (n *UserNotificationModel) toEntity() notificationEntity.UserNotification {

	notification := notificationEntity.NewUserNotificationFromDatabase(
		n.ID,
		n.UserID,
		n.NotificationTypeID,
		n.Data,
		n.CreatedAt,
		n.ViewedAt,
		n.UpdatedAt,
		n.Message,
		n.Title,
	)

	return notification
}

func NewNotificationRepository(sql *bun.DB, logger zerolog.Logger) *notificationPGRepository {
	return &notificationPGRepository{sql, logger}
}

func (r *notificationPGRepository) CreateUserNotification(ctx context.Context, userNotification notificationEntity.UserNotification) error {
	var dbNotification UserNotificationModel

	dbNotification = dbNotification.toDB(userNotification)
	_, err := r.db.NewInsert().Model(&dbNotification).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *notificationPGRepository) GetByUserID(ctx context.Context, userID string) ([]notificationEntity.UserNotification, error) {
	notificationModels := make([]UserNotificationModel, 0)
	err := r.db.NewSelect().
		Model(&notificationModels).
		Where("user_id = (?)", userID).
		OrderExpr("created_at DESC").
		Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("notificationPGRepository -> GetByUserID -> r.db.NewSelect(): %w", err)
	}

	userNotifications := make([]notificationEntity.UserNotification, 0, len(notificationModels))
	for _, notificationModel := range notificationModels {
		notificationEntity := notificationModel.toEntity()

		userNotifications = append(userNotifications, notificationEntity)
	}

	return userNotifications, nil
}

func (r *notificationPGRepository) GetByUserIDAndUserNotificationID(
	ctx context.Context,
	userID string,
	userNotificationID string,
) (notificationEntity.UserNotification, error) {
	var userNotificationModel UserNotificationModel
	err := r.db.NewSelect().
		Model(&userNotificationModel).
		Where("user_id = (?)", userID).
		Where("id = (?)", userNotificationID).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return notificationEntity.UserNotification{}, nil
		}
		return notificationEntity.UserNotification{}, fmt.Errorf("notificationPGRepository GetByUserIDAndUserNotificationID: %w", err)
	}

	return userNotificationModel.toEntity(), err
}

func (r *notificationPGRepository) MarkUserNotificationViewed(
	ctx context.Context,
	userID string,
	userNotificationID string,
) error {
	var userNotification UserNotificationModel
	userNotification.ViewedAt = time.Now()
	userNotification.UpdatedAt = time.Now()
	res, err := r.db.NewUpdate().
		Model(&userNotification).
		Column("viewed_at").
		Column("updated_at").
		Where("id = ?", userNotificationID).
		Where("user_id = ?", userID).
		Where("viewed_at is NULL").
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("notificationPGRepository MarkUserNotificationViewed -> r.db.NewUpdate: %w", err)
	}
	_, err = res.RowsAffected()
	return err
}

func (r *notificationPGRepository) MarkAllUserNotificationViewed(
	ctx context.Context,
	userID string,
) error {
	updates := map[string]interface{}{
		"viewed_at":  time.Now(),
		"updated_at": time.Now(),
	}
	_, err := r.db.NewUpdate().
		Model(&updates).
		Table("notifications").
		Where("user_id = (?)", userID).
		Where("viewed_at is NULL").
		Exec(ctx)

	return err
}

func (r *notificationPGRepository) DeleteUserNotification(
	ctx context.Context,
	userID string,
	userNotificationID string,
) error {
	var user UserNotificationModel
	_, err := r.db.NewDelete().Model(&user).Where("id = ?", userNotificationID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("notificationPGRepository DeleteUserNotification -> NewDelete: %w", err)
	}

	return err
}
