package repository

import (
	"context"
	"time"

	"github.com/PongponZ/test-container-golang/pkg/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TodoRepository interface {
	Create(entity.Todo) (string, error)
	Gets() ([]entity.Todo, error)
	Update(entity.Todo) error
	Delete(string) error
}

type todo struct {
	collection *mongo.Collection
}

func NewToDo(db *mongo.Database) TodoRepository {
	todoCollection := db.Collection("todo")

	return &todo{
		collection: todoCollection,
	}
}

func (t *todo) Create(task entity.Todo) (string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	result, err := t.collection.InsertOne(ctx, task)
	if err != nil {
		return "", err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()
	return id, err
}

func (t *todo) Gets() ([]entity.Todo, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	cursor, err := t.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var todos []entity.Todo
	if err = cursor.All(ctx, &todos); err != nil {
		return nil, err
	}

	return todos, nil
}

func (t *todo) Update(task entity.Todo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"_id": task.ID}
	update := bson.M{"$set": bson.M{"title": task.Title, "description": task.Description}}

	_, err := t.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (t *todo) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}

	_, err = t.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
