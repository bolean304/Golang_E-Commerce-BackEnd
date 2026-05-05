package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DatabaseName       = "Ecommerce"
	UsersCollection    = "Users"
	ProductsCollection = "Products"
)

func DBSet() *mongo.Client {
	client, err := NewMongoClient()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("successfully connected to mongodb")
	return client
}

var Client *mongo.Client = DBSet()

func MongoURI() string {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://admin:password@mongo:27017"
	}
	return mongoURI
}

func NewMongoClient() (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoURI()))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = client.Connect(ctx); err != nil {
		return nil, err
	}
	if err = client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, err
	}
	return client, nil
}

func UserData(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database(DatabaseName).Collection(collectionName)
	return collection
}
func ProductData(client *mongo.Client, collectionName string) *mongo.Collection {
	var productCollection *mongo.Collection = client.Database(DatabaseName).Collection(collectionName)
	return productCollection
}
