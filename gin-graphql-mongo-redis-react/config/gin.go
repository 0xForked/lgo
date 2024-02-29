package config

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	lmgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	lsredis "github.com/ulule/limiter/v3/drivers/store/redis"
	"go.uber.org/zap"
)

var allowOrigins = []string{
	"http://localhost:3000",
	"http://localhost:8000",
}

var allowHeaders = []string{
	"Content-Type",
	"Content-Length",
	"Accept-Encoding",
	"Authorization",
	"Cache-Control",
	"Origin",
	"Cookie",
}

func (c *Config) InitGinEngine() *Config {
	serverEngineSingleton.Do(func() {
		log.Printf("Trying to init engine (GIN %s) . . . .", gin.Version)
		// set gin mode
		gin.SetMode(func() string {
			if c.AppEnvironment == "local" {
				return gin.DebugMode
			}
			return gin.ReleaseMode
		}())
		// set global variables
		GinEngine = gin.Default()
		// set cors middleware
		GinEngine.Use(cors.New(cors.Config{
			AllowOrigins:     allowOrigins,
			AllowMethods:     []string{"GET, POST, PATCH, DELETE"},
			AllowHeaders:     allowHeaders,
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				return origin == "http://localhost:3000"
			},
			MaxAge: 12 * time.Hour,
		}))
		if c.AppEnvironment != "local" {
			// setup sentry middleware
			// setup rate limiter
			rate, err := limiter.NewRateFromFormatted("100-M")
			if err != nil {
				log.Fatalf("RATELIMIT_ERROR: %s", err.Error())
			}
			store, err := lsredis.NewStoreWithOptions(RedisPool,
				limiter.StoreOptions{Prefix: "social_app"})
			if err != nil {
				log.Fatalf("RATELIMIT_ERROR: %s", err.Error())
			}
			GinEngine.ForwardedByClientIP = true
			GinEngine.Use(lmgin.NewMiddleware(limiter.New(store, rate)))
			// setup logger
			logger, err := zap.NewProduction()
			if err != nil {
				log.Fatalf("LOGGER_ERROR: %s", err.Error())
			}
			defer func() { _ = logger.Sync() }()
			GinEngine.Use(ginzap.Ginzap(logger, time.RFC3339, true))
			GinEngine.Use(ginzap.RecoveryWithZap(logger, true))
		}
	})
	return c
}
