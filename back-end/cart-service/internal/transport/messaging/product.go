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
	productCreatedDurableConsumerName = "cart-service-product-created"
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
}

type ProductCreatedEvent struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (d *productMessagingHandlers) ProductCreatedListener() {
	d.logger.Info().Msg("ProductListener initialized")
	handler := func(n *nats.Msg) error {
		messageData := n.Data
		log.Info().Msg("ProductListener -> Received a message: " + string(messageData))

		var productEvent ProductCreatedEvent
		err := json.Unmarshal(messageData, &productEvent)
		if err != nil {
			log.Error().Msg("ProductListener -> Error in unmarshalling the message")
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
			log.Error().Err(err).Msg("ProductListener -> d.appService.CreateProduct")
			return err
		}
		return nil
	}
	d.natsClient.SubscribeDurable(productCreatedSubject, productStream, productCreatedDurableConsumerName, handler)
}
