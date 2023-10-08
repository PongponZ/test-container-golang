//go:build integration
// +build integration

package usecase_test

import (
	"context"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoContainer struct {
	testcontainers.Container
	URI string
}

// there are two ways to create test container
// 1. create from container request
func setupMongodbManual(ctx context.Context) (*mongoContainer, error) {
	//create request for start container
	req := testcontainers.ContainerRequest{
		Image:        "mongo:4.4.3",
		ExposedPorts: []string{"27017/tcp"},
		Env: map[string]string{
			"MONGO_INITDB_DATABASE":      "test-container",
			"MONGO_INITDB_ROOT_USERNAME": "root",
			"MONGO_INITDB_ROOT_PASSWORD": "root",
		},
		WaitingFor: wait.ForExec([]string{"echo", "db.runCommand('ping').ok", "|", "mongosh", "localhost:27017/test", "--quiet"}), // for check container is ready
	}

	//start container similar to the docker run command.
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	ip, err := container.Host(ctx) // get host of container
	if err != nil {
		return nil, err
	}

	mapPort, err := container.MappedPort(ctx, "27017") // get mapped port of container
	if err != nil {
		return nil, err
	}

	//create uri for connect to mongodb
	uri := "mongodb://" + ip + ":" + mapPort.Port()

	return &mongoContainer{
		Container: container,
		URI:       uri,
	}, nil
}

// 2. create from build-in module
func setMongodbContainer(ctx context.Context) (*mongoContainer, error) {
	customReq := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "mongo:4",
			Env: map[string]string{
				"MONGO_INITDB_DATABASE":      "test-container",
				"MONGO_INITDB_ROOT_USERNAME": "root",
				"MONGO_INITDB_ROOT_PASSWORD": "root",
			},
		},
		Started: true,
	}

	container, err := mongodb.RunContainer(ctx, testcontainers.CustomizeRequest(customReq))
	if err != nil {
		return nil, err
	}

	//create uri for connect to mongodb
	uri, err := container.ConnectionString(ctx)
	if err != nil {
		return nil, err
	}

	return &mongoContainer{
		Container: container,
		URI:       uri,
	}, nil
}

func connectMongodb(ctx context.Context, uri string) (*mongo.Client, error) {
	credential := options.Credential{
		Username: "root",
		Password: "root",
	}
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetAuth(credential))
	if err != nil {
		panic(err)
	}
	return mongoClient, nil
}

func TestTodoUsecase_Integration(t *testing.T) {
	ctx := context.Background()

	//setup mongodb container
	container, err := setMongodbContainer(ctx)
	if err != nil {
		panic(err)
	}

	//mongo client
	_, err = connectMongodb(ctx, container.URI)
	if err != nil {
		panic(err)
	}

	//clean up
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	t.Run("Create", func(t *testing.T) {

	})

}
