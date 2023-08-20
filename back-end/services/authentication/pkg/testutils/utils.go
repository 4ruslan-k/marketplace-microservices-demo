package testutils

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type container struct {
	testcontainers.Container
	URI string
}

type image struct {
	name string
	port string
}

var MongoDB = image{
	name: "mongo:6",
	port: "27017",
}

func InitializeMongoContainer(ctx context.Context) (testcontainers.Container, string, error) {
	var mongoImage = image{
		name: "mongo:6",
		port: "27017",
	}
	start := time.Now()
	req := testcontainers.ContainerRequest{
		Image:        mongoImage.name,
		ExposedPorts: []string{mongoImage.port + "/tcp"},
		WaitingFor:   wait.ForListeningPort(nat.Port(mongoImage.port)),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", err
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, "", err
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(mongoImage.port))
	if err != nil {
		return nil, "", err
	}

	uri := fmt.Sprintf("mongodb://%s:%s", hostIP, mappedPort.Port())

	_, err = mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, "", err
	}
	elapsed := time.Since(start)
	log.Printf("TestContainers: container %s is now running at %s\n", req.Image, uri)
	log.Printf("TestContainers: container took %s\n", elapsed)
	return container, uri, nil
}
