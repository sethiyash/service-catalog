package db

import (
	"context"
	"log"
	"service-catalog/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupDatabase() (*mongo.Client, error) {
	config.LoadConfig()

	mongoURI := config.GetEnv("MONGODB_URI", "")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable is not set")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")

	return client, nil
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database("service_catalog").Collection(collectionName)
}
