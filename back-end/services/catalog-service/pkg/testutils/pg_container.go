package testutils

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type testDatabase struct {
	instance testcontainers.Container
}

func (db *testDatabase) port() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	p, err := db.instance.MappedPort(ctx, "5432")
	if err != nil {
		return 0, err
	}
	return p.Int(), nil
}

type image struct {
	name string
	port string
}

var pgImage = image{
	name: "postgres:12",
	port: "5432",
}

func InitializePGContainer(ctx context.Context) (testcontainers.Container, string, error) {
	req := testcontainers.ContainerRequest{
		Image:        pgImage.name,
		ExposedPorts: []string{pgImage.port + "/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "postgres",
		},
		WaitingFor: wait.ForListeningPort(nat.Port(pgImage.port)),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return nil, "", fmt.Errorf("InitializePGContainer -> testcontainers.GenericContainer: %w", err)
	}

	testDB := testDatabase{instance: container}
	port, err := testDB.port()

	if err != nil {
		return nil, "", err
	}
	uri := fmt.Sprintf("postgres://postgres:postgres@127.0.0.1:%d/postgres?sslmode=disable", port)
	return container, uri, nil
}
