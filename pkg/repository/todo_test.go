//go:build integration
// +build integration

package repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/PongponZ/test-container-golang/pkg/entity"
	"github.com/PongponZ/test-container-golang/pkg/repository"
	"github.com/stretchr/testify/assert"
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
func setupMongoDB(ctx context.Context) (*mongoContainer, error) {
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
func setupMongoDBWithModule(ctx context.Context) (*mongoContainer, error) {
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

	fmt.Println("connection string: ", uri)

	return &mongoContainer{
		Container: container,
		URI:       uri,
	}, nil
}

func TestTodoRepository(t *testing.T) {
	ctx := context.Background()

	// mongoContainer, err := setupMongoDB(ctx)
	// if err != nil {
	// 	panic(err)
	// }

	mongoContainer, err := setupMongoDBWithModule(ctx)
	if err != nil {
		panic(err)
	}

	//connect to mongodb
	credential := options.Credential{
		Username: "root",
		Password: "root",
	}
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoContainer.URI).SetAuth(credential))
	if err != nil {
		panic(err)
	}

	cleanup := func() {
		if err := mongoContainer.Container.Terminate(ctx); err != nil {
			panic(err)
		}
	}

	t.Run("Create", func(t *testing.T) {
		var todoRepository repository.TodoRepository

		berforeEach := func() {
			todoRepository = repository.NewToDo(mongoClient)
		}

		t.Run("Should be able to insert document to mongodb", func(t *testing.T) {
			berforeEach()

			result, err := todoRepository.Create(entity.Todo{
				ID:          1,
				Title:       "my title",
				Description: "my description",
			})

			assert.NotNil(t, result)
			assert.NoError(t, err)
		})

		t.Run("Should be able to update document to mongodb", func(t *testing.T) {
			berforeEach()

			result, err := todoRepository.Create(entity.Todo{
				ID:          1,
				Title:       "my title",
				Description: "my description",
			})

			assert.NotNil(t, result)
			assert.NoError(t, err)
		})
	})

	t.Cleanup(cleanup)
}
