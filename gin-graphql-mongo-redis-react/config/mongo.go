package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoMinPool  = 10
	mongoMaxPool  = 100
	mongoIdleTime = 60
)

func (c *Config) InitMongoDBConn() *Config {
	mongoSingleton.Do(func() {
		log.Println("Init database connection pool . . .")
		var client *mongo.Client
		var err error
		opts := options.Client().
			ApplyURI(c.MongoDSN).
			SetMinPoolSize(uint64(mongoMinPool)).
			SetMaxPoolSize(uint64(mongoMaxPool)).
			SetMaxConnIdleTime(time.Duration(mongoIdleTime) * time.Second)
		if client, err = mongo.Connect(context.Background(), opts); err != nil {
			log.Fatalf("MONGO_ERROR: %s", err.Error())
		}
		MongoDbPool = client.Database(c.MongoDBName)
		if err := MongoDbPool.RunCommand(context.TODO(), bson.D{
			{"ping", 1},
		}).Err(); err != nil {
			log.Fatalf("MONGO_ERROR: %s", err.Error())
		}
		log.Println("MongoDB conn pool ready")
	})
	return c
}
