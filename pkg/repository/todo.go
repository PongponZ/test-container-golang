package repository

import (
	"context"
	"time"

	"github.com/PongponZ/test-container-golang/pkg/entity"
	"go.mongodb.org/mongo-driver/mongo"
)

type TodoRepository interface {
	Create(entity.Todo) (*mongo.InsertOneResult, error)
}

type todo struct {
	mongoClient *mongo.Client
}

func NewToDo(mongodb *mongo.Client) TodoRepository {

	return &todo{
		mongoClient: mongodb,
	}
}

func (t *todo) Create(task entity.Todo) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	database := t.mongoClient.Database("test-container")
	todoCollection := database.Collection("todo")

	result, err := todoCollection.InsertOne(ctx, task)

	return result, err
}
