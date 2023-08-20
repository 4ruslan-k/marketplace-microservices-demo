package messaging

import (
	"context"
	"encoding/json"
	natsClient "shared/messaging/nats"

	notification "notification_service/internal/domain/entities/notification"
	applicationServices "notification_service/internal/services"
	socketServer "notification_service/internal/transport/http/socketio"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	userNotificationCreationSubject       = "notifications.created"
	refetchUserNotificationSubject        = "notifications.refetch"
	notificationCreateDurableConsumerName = "notification-service-notification-create"
	notificationStream                    = "notifications"
)

type NotificationMessagingHandlers interface {
	NotificationListener()
	Init()
}

type UserNotificationRefetchEvent struct {
	UserID string `json:"id"`
}

type notificationMessagingHandlers struct {
	natsClient   natsClient.NatsClient
	logger       zerolog.Logger
	appService   applicationServices.NotificationApplicationService
	socketServer *socketServer.SocketIOServer
}

func NewNotificationMessagingHandlers(
	natsClient natsClient.NatsClient,
	appService applicationServices.NotificationApplicationService,
	logger zerolog.Logger,
	socketServer *socketServer.SocketIOServer,
) *notificationMessagingHandlers {
	d := notificationMessagingHandlers{natsClient: natsClient, appService: appService, logger: logger, socketServer: socketServer}
	return &d
}

func (d *notificationMessagingHandlers) Init() {
	d.logger.Info().Msg("initializing NotificationMessagingHandlers")

	err := d.natsClient.CreateStream(notificationStream, "notifications.*")
	if err != nil {
		log.Error().Err(err).Msg("notificationMessagingHandlers Init -> u.natsClient.CreateStream")
	}

	d.NotificationListener()
	d.RefetchNotificationListener()
}

type UserNotificationCreatedEvent struct {
	NotificationTypeID string      `json:"notificationTypeID"`
	UserID             string      `json:"userID"`
	Data               interface{} `json:"data,omitempty"`
}

func (d *notificationMessagingHandlers) NotificationListener() {
	d.logger.Info().Msg("NotificationListener initialized")
	handler := func(n *nats.Msg) error {
		messageData := n.Data
		log.Info().Msg("NotificationListener -> Received a message: " + string(messageData))

		var notificationEvent UserNotificationCreatedEvent
		err := json.Unmarshal(messageData, &notificationEvent)
		if err != nil {
			log.Error().Msg("NotificationListener -> Error in unmarshalling the message")
			return err
		}
		err = d.appService.CreateUserNotification(context.Background(), notification.CreateUserNotificationParams{
			UserID:             notificationEvent.UserID,
			NotificationTypeID: notificationEvent.NotificationTypeID,
			Data:               notificationEvent.Data,
		})
		if err != nil {
			log.Error().Err(err).Msg("NotificationListener -> d.appService.CreateUserNotification")
			return err
		}
		return nil
	}
	d.natsClient.SubscribeDurable(userNotificationCreationSubject, notificationStream, notificationCreateDurableConsumerName, handler)
}

func (d *notificationMessagingHandlers) RefetchNotificationListener() {
	d.logger.Info().Msg("RefetchNotificationListener initialized")
	handler := func(n *nats.Msg) error {
		messageData := n.Data
		var notificationEvent UserNotificationRefetchEvent
		err := json.Unmarshal(messageData, &notificationEvent)
		if err != nil {
			log.Error().Err(err).Msg("NotificationListener -> d.appService.RefetchNotificationListener")
			return err
		}

		log.Info().Msg("RefetchNotificationListener -> Received a message: " + string(messageData))

		d.socketServer.SendEvent(notificationEvent.UserID)

		return nil
	}
	d.natsClient.SubscribeEphemeral(refetchUserNotificationSubject, handler)
}
