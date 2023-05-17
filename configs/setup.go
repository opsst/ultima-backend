package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(EnvMongoURI()))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Ultima MongoDB")
	return client
}

func ConnectDB2() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(EnvMongoURI2()))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Florxy MongoDB")
	return client
}

// Client instance
var DB *mongo.Client = ConnectDB()
var DB2 *mongo.Client = ConnectDB2()

// getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("Ultima-API").Collection(collectionName)
	return collection
}

func GetCollection2(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("Florxy").Collection(collectionName)
	return collection
}
