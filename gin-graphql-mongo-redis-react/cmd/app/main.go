package main

import (
	"github.com/spf13/viper"
	"learn.go/gin-graphql-mongo-redis/config"
	"learn.go/gin-graphql-mongo-redis/server"
)

func main() {
	viper.SetConfigFile(".env")
	config.LoadEnv().
		InitMongoDBConn().
		InitRedisConnection().
		InitGinEngine()
	server.Run()
}
