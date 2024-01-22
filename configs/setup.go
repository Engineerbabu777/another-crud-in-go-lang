package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/another-crud-in-go-lang/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
    client, err := mongo.NewClient(options.Client().ApplyURI(configs.EnvMongoUri()));

	if err!= nil {
		log.Fatal(err);
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second);

	err = client.Connect(ctx);

	if err != nil {
        log.Fatal(err)
    }

	ping := client.Ping(ctx, nil);

	if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Connected to MongoDB");

    return client

}


var DBClient *mongo.Client = ConnectDB(); 

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection{
   collection := client.Database("goLangApi").Collection(collectionName);
   return collection;
}