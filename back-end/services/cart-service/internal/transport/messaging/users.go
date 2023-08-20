package messaging

import (
	"context"
	"encoding/json"
	natsClient "shared/messaging/nats"

	"cart_service/internal/application/dto"
	applicationServices "cart_service/internal/application/services"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	customerEntity "cart_service/internal/domain/entities/customer"
)

const (
	userCreationSubject           = "users.created"
	userUpdateSubject             = "users.updated"
	userDeletionSubject           = "users.deleted"
	usersStreamName               = "users"
	userCreateDurableConsumerName = "cart-service-user-create"
	userDeleteDurableConsumerName = "cart-service-user-delete"
	userUpdateDurableConsumerName = "cart-service-user-update"
)

type UserMessagingHandlers interface {
	UserCreationListener()
	UserUpdateListener()
	UserDeletedListener()
	Init()
}

type userMessagingHandlers struct {
	natsClient natsClient.NatsClient
	logger     zerolog.Logger
	appService applicationServices.CustomerApplicationService
}

func NewCustomerMessagingHandlers(
	natsClient natsClient.NatsClient,
	appService applicationServices.CustomerApplicationService,
	logger zerolog.Logger,
) *userMessagingHandlers {
	u := userMessagingHandlers{natsClient: natsClient, appService: appService, logger: logger}
	return &u
}

func (u *userMessagingHandlers) Init() {
	u.logger.Info().Msg("initializing UserMessagingHandlers")
	err := u.natsClient.CreateStream(usersStreamName, "users.*")
	if err != nil {
		log.Error().Err(err).Msg("Init -> u.natsClient.CreateStream")
	}

	u.UserCreationListener()
	u.UserUpdateListener()
	u.UserDeletedListener()
}

type UserCreatedEvent struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	ID    string `json:"id"`
}

type UserUpdatedEvent struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type UserDeletedEvent struct {
	ID string `json:"id"`
}

func (u *userMessagingHandlers) UserCreationListener() {
	u.logger.Info().Msg("UserCreationListener initialized")
	handler := func(n *nats.Msg) error {
		messageData := n.Data
		log.Info().Msg("UserCreationListener -> Received a message: " + string(messageData))

		var userCreatedEvent UserCreatedEvent
		err := json.Unmarshal(messageData, &userCreatedEvent)
		if err != nil {
			log.Error().Msg("UserCreationListener -> Error in unmarshalling the message")
			return err
		}
		_, err = u.appService.CreateCustomer(context.Background(), customerEntity.CreateCustomerParams{
			Name:  userCreatedEvent.Name,
			Email: userCreatedEvent.Email,
			ID:    userCreatedEvent.ID,
		})
		if err != nil {
			log.Error().Msg("UserCreationListener -> u.appService.CreateCustomer")
			return err
		}
		return nil
	}
	u.natsClient.SubscribeDurable(userCreationSubject, usersStreamName, userCreateDurableConsumerName, handler)
}

func (u *userMessagingHandlers) UserUpdateListener() {
	u.logger.Info().Msg("UserUpdateListener initialized")
	handler := func(n *nats.Msg) error {
		messageData := n.Data
		log.Info().Msg("UserUpdateListener -> Received a message: " + string(messageData))

		var userUpdatedEvent UserUpdatedEvent
		err := json.Unmarshal(messageData, &userUpdatedEvent)
		if err != nil {
			log.Error().Msg("UserUpdateListener -> Error in unmarshalling the message")
			return err
		}
		err = u.appService.UpdateCustomer(
			context.Background(),
			dto.UpdateCustomerInput{
				Name: userUpdatedEvent.Name,
				ID:   userUpdatedEvent.ID,
			},
		)
		if err != nil {
			log.Error().Err(err).Msg("UserUpdateListener ->u.appService.UpdateCustomer")
			return err
		}
		return nil
	}
	u.natsClient.SubscribeDurable(userUpdateSubject, usersStreamName, userUpdateDurableConsumerName, handler)
}

func (u *userMessagingHandlers) UserDeletedListener() {
	u.logger.Info().Msg("UserDeletedListener initialized")
	handler := func(n *nats.Msg) error {
		messageData := n.Data
		log.Info().Msg("UserDeletedListener -> Received a message: " + string(messageData))

		var userDeletedEvent UserDeletedEvent
		err := json.Unmarshal(messageData, &userDeletedEvent)
		if err != nil {
			log.Error().Msg("UserUpdateListener -> Error in unmarshalling the message")
			return err
		}
		u.appService.DeleteCustomer(
			context.Background(),
			dto.DeleteCustomerInput{
				ID: userDeletedEvent.ID,
			},
		)
		return nil
	}
	u.natsClient.SubscribeDurable(userDeletionSubject, usersStreamName, userDeleteDurableConsumerName, handler)
}
