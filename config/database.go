package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Db *mongo.Database

func InitDB(conf Configuration) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	credential := options.Credential{
		Username: conf.MONGO_USERNAME,
		Password: conf.MONGO_PASSWORD,
	}
	clientOptions := options.Client().ApplyURI(conf.MONGO_URL).SetAuth(credential)

	mongoClient, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Println(err, "Connection in mongodb")
		log.Println(mongoClient.Ping(ctx, readpref.Primary()), "Ping error in mongoconnect")
	}

	Db = mongoClient.Database(conf.MONGO_DB_NAME)
}
