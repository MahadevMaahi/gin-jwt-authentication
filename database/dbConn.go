package database

import (
	"fmt"
	"log"
	"time"
	"os"
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DbClient struct {}

type DbKit interface {
	OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection
}

func DBInstance() *mongo.Client {
	fmt.Println("Connecting to DB")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error Loading .env file")
	}

	MongoDb := os.Getenv("MONGODB_URL")

	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err1 := client.Connect(ctx)
	if err1 != nil {
		log.Fatal(err1)
	}

	fmt.Println("Connected to the Data Base")

	return client
}

var Client *mongo.Client = DBInstance()

func (dbClient *DbClient) OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	fmt.Println("Openning the Collection")
	var collection *mongo.Collection = client.Database("cluster0").Collection(collectionName)
	return collection
}