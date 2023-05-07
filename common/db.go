package common

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func GetDbCollection(col string) *mongo.Collection {
	return db.Collection(col)
}

func InitDB() error {
	uri := os.Getenv("MONGODB_URI")

	if uri == "" {
		log.Fatal("You must set your 'MONGODB URI' env variables")
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))

	if err != nil {
		return err
	}

	db = client.Database("demo")

	return nil
}

func CloseDB() error {
	return db.Client().Disconnect(context.Background())
}
