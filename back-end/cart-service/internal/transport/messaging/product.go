package messaging

import (
	natsClient "cart_service/pkg/messaging/nats"
	"context"
	"encoding/json"
	"time"

	applicationServices "cart_service/internal/application/services"
	productEntity "cart_service/internal/domain/entities/product"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	productCreatedSubject             = "products.created"
	productDeletedSubject             = "products.deleted"
	productUpdatedSubject             = "products.updated"
	productCreatedDurableConsumerName = "cart-service-product-created"
	productDeletedDurableConsumerName = "cart-service-product-deleted"
	productUpdatedDurableConsumerName = "cart-service-product-updated"
	productStream                     = "products"
)

type ProductMessagingHandlers interface {
	ProductCreatedListener()
	Init()
}

var _ ProductMessagingHandlers = (*productMessagingHandlers)(nil)

type productMessagingHandlers struct {
	natsClient natsClient.NatsClient
	logger     zerolog.Logger
	appService applicationServices.ProductApplicationService
}

func NewProductMessagingHandlers(
	natsClient natsClient.NatsClient,
	appService applicationServices.ProductApplicationService,
	logger zerolog.Logger,
) *productMessagingHandlers {
	d := productMessagingHandlers{natsClient: natsClient, appService: appService, logger: logger}
	return &d
}

func (d *productMessagingHandlers) Init() {
	d.logger.Info().Msg("initializing ProductMessagingHandlers")

	err := d.natsClient.CreateStream(productStream, "products.*")
	if err != nil {
		log.Error().Err(err).Msg("productMessagingHandlers Init -> u.natsClient.CreateStream")
	}

	d.ProductCreatedListener()
	d.ProductDeletedListener()
	d.ProductUpdatedListener()
}

type ProductCreatedEvent struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductDeletedEvent struct {
	ID string `json:"id"`
}

type ProductUpdatedEvent struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (d *productMessagingHandlers) ProductCreatedListener() {
	d.logger.Info().Msg("ProductCreatedListener initialized")
	handler := func(n *nats.Msg) error {
		messageData := n.Data
		log.Info().Msg("ProductCreatedListener -> Received a message: " + string(messageData))

		var productEvent ProductCreatedEvent
		err := json.Unmarshal(messageData, &productEvent)
		if err != nil {
			log.Error().Msg("ProductCreatedListener -> Error in unmarshalling the message")
			return err
		}
		err = d.appService.CreateProduct(context.Background(), productEntity.CreateProductParams{
			ID:        productEvent.ID,
			Name:      productEvent.Name,
			Price:     productEvent.Price,
			Quantity:  productEvent.Quantity,
			CreatedAt: productEvent.CreatedAt,
			UpdatedAt: productEvent.UpdatedAt,
		})
		if err != nil {
			log.Error().Err(err).Msg("ProductCreatedListener -> d.appService.CreateProduct")
			return err
		}
		return nil
	}
	d.natsClient.SubscribeDurable(productCreatedSubject, productStream, productCreatedDurableConsumerName, handler)
}

func (d *productMessagingHandlers) ProductDeletedListener() {
	d.logger.Info().Msg("ProductDeletedListener initialized")
	handler := func(n *nats.Msg) error {
		messageData := n.Data
		log.Info().Msg("ProductDeletedListener -> Received a message: " + string(messageData))

		var productDeletedEvent ProductDeletedEvent
		err := json.Unmarshal(messageData, &productDeletedEvent)
		if err != nil {
			log.Error().Msg("ProductDeletedListener -> Error in unmarshalling the message")
			return err
		}
		err = d.appService.DeleteProduct(context.Background(), productDeletedEvent.ID)
		if err != nil {
			log.Error().Err(err).Msg("ProductDeletedListener -> d.appService.CreateProduct")
			return err
		}
		return nil
	}
	d.natsClient.SubscribeDurable(productDeletedSubject, productStream, productDeletedDurableConsumerName, handler)
}

func (d *productMessagingHandlers) ProductUpdatedListener() {
	d.logger.Info().Msg("ProductUpdatedListener initialized")
	handler := func(n *nats.Msg) error {
		messageData := n.Data
		log.Info().Msg("ProductUpdatedListener -> Received a message: " + string(messageData))

		var productUpdatedEvent ProductUpdatedEvent
		err := json.Unmarshal(messageData, &productUpdatedEvent)
		if err != nil {
			log.Error().Msg("ProductUpdatedListener -> Error in unmarshalling the message")
			return err
		}
		err = d.appService.UpdateProduct(context.Background(), productEntity.CreateProductParams{
			ID:        productUpdatedEvent.ID,
			Name:      productUpdatedEvent.Name,
			Price:     productUpdatedEvent.Price,
			Quantity:  productUpdatedEvent.Quantity,
			CreatedAt: productUpdatedEvent.CreatedAt,
			UpdatedAt: productUpdatedEvent.UpdatedAt,
		})
		if err != nil {
			log.Error().Err(err).Msg("ProductUpdatedListener -> d.appService.CreateProduct")
			return err
		}
		return nil
	}
	d.natsClient.SubscribeDurable(productUpdatedSubject, productStream, productUpdatedDurableConsumerName, handler)
}
