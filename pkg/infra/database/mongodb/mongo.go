package mongodb

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoAdapter interface {
	Disconnect(ctx context.Context)
	GetDatabase(db string) *mongo.Database
}

type mongodb struct {
	client *mongo.Client
	uri    string
}

func NewConnect(ctx context.Context, uri string) MongoAdapter {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("could not connect to mongo: %v\n", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("could not connect to mongo: %v\n", err)
	}

	return &mongodb{
		client: client,
		uri:    uri,
	}
}

func (m *mongodb) Disconnect(ctx context.Context) {
	if err := m.client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

func (m *mongodb) GetDatabase(db string) *mongo.Database {
	return m.client.Database(db)
}
