package storage

import (
	"authentication/config"
	"context"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient(logger zerolog.Logger, config *config.Config) *mongo.Database {
	opts := options.Client()
	opts.ApplyURI(config.MongoURI)
	mongoClient, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		logger.Panic().Err(err).Msg("initiating mongodb client")
	}
	mongoDB := mongoClient.Database(config.MongoDatabaseName)
	return mongoDB
}
