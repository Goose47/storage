package db

import (
	"Goose47/storage/config"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

var Client *mongo.Client

func GetCollection() *mongo.Collection {
	return Client.Database(config.DBConfig.DBName).Collection(config.DBConfig.DBColl)
}

func Init() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var err error
	Client, err = mongo.Connect(ctx, options.Client().ApplyURI(config.DBConfig.Url))
	if err != nil {
		log.Fatal(err)
	}

	err = Client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
}
