package nats

import (
	nats "shared/messaging/nats"
)

type NatsClient interface {
	nats.NatsClient
}
