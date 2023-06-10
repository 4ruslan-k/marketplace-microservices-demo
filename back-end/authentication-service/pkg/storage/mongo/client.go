package storage

import (
	"authentication_service/config"
	"context"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

func NewMongoClient(logger zerolog.Logger, config *config.Config) *mongo.Database {
	opts := options.Client()
	opts.Monitor = otelmongo.NewMonitor()
	opts.ApplyURI(config.MongoURI)
	mongoClient, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		logger.Panic().Err(err).Msg("initiating mongodb client")
	}
	mongoDB := mongoClient.Database(config.MongoDatabaseName)
	return mongoDB
}
