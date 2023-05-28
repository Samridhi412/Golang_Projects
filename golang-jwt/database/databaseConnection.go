package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client{
	err := godotenv.Load(".env")
	if err != nil{
		log.Fatal("Error loading .env file")
	}
	//if database doesnt exist mongo creates for you
	MongoDb := os.Getenv("MONGODB_URL")
	// Create a new MongoDB client object
	client, err:=mongo.NewClient(options.Client().ApplyURI(MongoDb))
	if err != nil{
		log.Fatal(err)
	}
	//gin.context- It provides methods for accessing information about the request and manipulating the response.
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println("Connected to mongodb")
	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection{
	//collection: A collection is a grouping of MongoDB documents that are stored in a particular database
	var collection *mongo.Collection = client.Database("c").Collection(collectionName)
	return collection
}