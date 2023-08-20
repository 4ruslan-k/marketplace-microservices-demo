package nats

import (
	"errors"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type NatsClient interface {
	PublishMessage(subject string, message string)
	PublishMessageEphemeral(subject string, message string)
	CreateStream(streamName string, streamSubjects string) error
	SubscribeDurable(subject string, streamName string, consumerName string, handler func(m *nats.Msg) error)
	AddConsumer(streamName string, consumerName string, subject string) error
	SubscribeEphemeral(subject string, handler func(m *nats.Msg) error)
}

type natsClient struct {
	nc     *nats.Conn
	js     nats.JetStreamContext
	logger zerolog.Logger
}

func (n natsClient) PublishMessage(subject, message string) {
	n.logger.Debug().Msgf("Publishing message to %s", subject)
	_, err := n.js.Publish(subject, []byte(message))

	if err != nil {
		log.Error().Err(err).Msg("createTask -> js.Publish")
	}
}

func (n natsClient) PublishMessageEphemeral(subject, message string) {
	n.logger.Debug().Msgf("Publishing ephemeral message to %s", subject)
	err := n.nc.Publish(subject, []byte(message))

	if err != nil {
		log.Error().Err(err).Msg("createTask -> js.Publish")
	}
}

func newNatsConnection(uri string) *nats.Conn {
	var nc *nats.Conn
	var err error

	for i := 0; i < 5; i++ {
		nc, err = nats.Connect(uri)
		if err == nil {
			break
		}
		log.Error().Err(err).Msg("Error establishing connection to NATS")
		log.Info().Msgf("Waiting before connecting to NATS at: %s", uri)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		panic(err)
	}
	log.Info().Msgf("Successfully connected to NATS at: %s", uri)
	return nc
}

func NewNatsClient() *natsClient {
	uri := os.Getenv("NATS_URI")
	nc := newNatsConnection(uri)
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		panic(err)
	}
	serverNats := natsClient{nc: nc, js: js}
	return &serverNats
}

func (n natsClient) CreateStream(streamName string, streamSubjects string) error {
	stream, err := n.js.StreamInfo(streamName)
	if err != nil {
		if errors.Is(nats.ErrStreamNotFound, err) {
			// stream not found, create it
			if stream == nil {
				n.logger.Debug().Msgf("CreateStream -> Creating stream: %s", streamName)
				_, err = n.js.AddStream(&nats.StreamConfig{
					Name:      streamName,
					Subjects:  []string{streamSubjects},
					Retention: nats.InterestPolicy,
				})
				if err != nil {
					return err
				}
			}
			return nil
		}
		n.logger.Err(err).Msg("CreateStream -> n.js.StreamInfo")
		return err
	}
	return nil
}

func (n natsClient) AddConsumer(streamName string, consumerName string, subject string) error {
	_, err := n.js.UpdateConsumer(streamName, &nats.ConsumerConfig{
		Durable:       consumerName,
		AckPolicy:     nats.AckExplicitPolicy,
		FilterSubject: subject,
	})
	return err
}

func (n natsClient) SubscribeDurable(subject string, streamName string, consumerName string, handler func(m *nats.Msg) error) {
	log.Logger.Info().
		Str("subject", subject).
		Str("streamName", streamName).
		Str("consumerName", consumerName).
		Msg("Subscribing to subject")
	_, err := n.js.Subscribe(subject, func(m *nats.Msg) {
		err := m.Ack()

		if err != nil {
			n.logger.Error().Err(err).Msg("Error acking message")
		}

		handler(m)
	}, nats.BindStream(streamName), nats.Durable(consumerName), nats.AckExplicit())

	if err != nil {
		n.logger.Err(err).Msg("natsClient -> Subscribe")
		return
	}
}

func (n natsClient) SubscribeEphemeral(subject string, handler func(m *nats.Msg) error) {
	log.Logger.Info().
		Str("subject", subject).
		Msg("Subscribing to subject")
	_, err := n.nc.Subscribe(subject, func(m *nats.Msg) {
		err := m.Ack()

		if err != nil {
			n.logger.Error().Err(err).Msg("Error acking message")
		}

		handler(m)
	})

	if err != nil {
		n.logger.Err(err).Msg("natsClient -> SubscribeEphemeral")
		return
	}
}
