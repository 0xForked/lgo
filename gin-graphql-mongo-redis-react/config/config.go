package config

import (
	"errors"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	configSingleton, mongoSingleton, redisSingleton,
	serverEngineSingleton sync.Once

	Instance    *Config
	MongoDbPool *mongo.Database
	RedisPool   *redis.Client
	GinEngine   *gin.Engine
)

type Config struct {
	AppName        string `mapstructure:"APP_NAME"`
	AppVersion     string `mapstructure:"APP_VERSION"`
	AppEnvironment string `mapstructure:"APP_ENVIRONMENT"`
	ServerPort     string `mapstructure:"SERVER_PORT"`
	MongoDBName    string `mapstructure:"MONGO_DB_NAME"`
	MongoDSN       string `mapstructure:"MONGO_DSN_URL"`
	RedisDSN       string `mapstructure:"REDIS_DSN_URL"`
}

func LoadEnv() *Config {
	configSingleton.Do(func() {
		log.Println("Load configuration file . . .")
		// find environment file
		viper.AutomaticEnv()
		// error handling for specific case
		if err := viper.ReadInConfig(); err != nil {
			var configFileNotFoundError viper.ConfigFileNotFoundError
			if errors.As(err, &configFileNotFoundError) {
				// Config file not found; ignore error if desired
				log.Fatalln(".env file not found!, please copy .env.example and paste as .env")
			}
			log.Fatalf("ENV_ERROR: %s", err.Error())
		}
		// extract config to struct
		if err := viper.Unmarshal(&Instance); err != nil {
			log.Fatalf("ENV_ERROR: %s", err.Error())
		}
		// notify that config file is ready
		log.Println("Configuration file ready")
	})
	return Instance
}
