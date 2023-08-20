package controllers

import (
	"encoding/json"
	"notification_service/config"
	"notification_service/internal/domain/entities/notification"
	applicationServices "notification_service/internal/services"

	"github.com/rs/zerolog"

	"github.com/gin-gonic/gin"

	httpDto "notification_service/internal/transport/http/dto"
	httpErrors "shared/errors/http"
)

type NotificationControllers struct {
	ApplicationService applicationServices.NotificationApplicationService
	Logger             zerolog.Logger
	Config             *config.Config
}

func NewNotificationController(
	appService applicationServices.NotificationApplicationService,
	logger zerolog.Logger,
	config *config.Config,
) *NotificationControllers {
	return &NotificationControllers{
		ApplicationService: appService,
		Logger:             logger,
		Config:             config,
	}
}

type AuthInfo struct {
	UserID string `json:"user_id"`
}

func (r *NotificationControllers) GetNotificationsByNotificationID(c *gin.Context) {
	var authInfo AuthInfo
	authValue := c.Request.Header.Get("X-Authentication-Info")
	json.Unmarshal([]byte(authValue), &authInfo)

	notifications, err := r.ApplicationService.GetNotificationsByUserID(c.Request.Context(), authInfo.UserID)
	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}

	handleResponseWithBody(c, notifications)
}

func (r *NotificationControllers) CreateUserNotification(c *gin.Context) {
	var createUserNotificationInput httpDto.CreateUserNotificationInput
	if err := c.ShouldBindJSON(&createUserNotificationInput); err != nil {
		httpErrors.BadRequest(c, err.Error())
		return
	}

	err := r.ApplicationService.CreateUserNotification(c.Request.Context(), notification.CreateUserNotificationParams{
		UserID:             createUserNotificationInput.UserID,
		NotificationTypeID: createUserNotificationInput.NotificationTypeID,
		Data:               createUserNotificationInput.Data,
	})
	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}
	handleOkResponse(c)
}

func (r *NotificationControllers) ViewNotification(c *gin.Context) {

	notificationId, _ := c.Params.Get("notificationId")

	var authInfo AuthInfo
	authValue := c.Request.Header.Get("X-Authentication-Info")
	json.Unmarshal([]byte(authValue), &authInfo)

	err := r.ApplicationService.ViewNotification(c.Request.Context(), authInfo.UserID, notificationId)
	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}
	handleOkResponse(c)
}

func (r *NotificationControllers) DeleteUserNotification(c *gin.Context) {

	notificationId, _ := c.Params.Get("notificationId")

	var authInfo AuthInfo
	authValue := c.Request.Header.Get("X-Authentication-Info")
	json.Unmarshal([]byte(authValue), &authInfo)

	err := r.ApplicationService.DeleteUserNotification(c.Request.Context(), authInfo.UserID, notificationId)
	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}
	handleOkResponse(c)
}

func (r *NotificationControllers) ViewAllNotifications(c *gin.Context) {
	var authInfo AuthInfo
	authValue := c.Request.Header.Get("X-Authentication-Info")
	json.Unmarshal([]byte(authValue), &authInfo)

	err := r.ApplicationService.ViewAllNotifications(c.Request.Context(), authInfo.UserID)
	if err != nil {
		httpErrors.RespondWithError(c, err)
		return
	}
	handleOkResponse(c)
}
