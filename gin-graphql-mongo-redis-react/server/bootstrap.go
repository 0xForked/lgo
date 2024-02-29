package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	ginhealthcheck "github.com/tavsec/gin-healthcheck"
	healthcheck "github.com/tavsec/gin-healthcheck/checks"
	healthcheckconfig "github.com/tavsec/gin-healthcheck/config"
	"learn.go/gin-graphql-mongo-redis/config"
	"learn.go/gin-graphql-mongo-redis/server/common"
	"learn.go/gin-graphql-mongo-redis/server/graph"
)

var (
	CertFilePath = "./Cert/pokebook.server+3.pem"
	KeyFilePath  = "./Cert/pokebook.server+3-key.pem"
)

func Run() {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	// router engine
	routerEngine := config.GinEngine
	// register public routes
	registerRoutes(routerEngine, ctx)
	// load tls certificates
	serverTLSCert, err := tls.LoadX509KeyPair(CertFilePath, KeyFilePath)
	if err != nil {
		log.Fatalf("Error loading certificate and key file: %v", err)
	}
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{serverTLSCert}}
	// server defines parameters for running an HTTP server.
	server := &http.Server{
		Addr:              config.Instance.ServerPort,
		Handler:           routerEngine,
		ReadHeaderTimeout: time.Second * 10,
		TLSConfig:         tlsConfig,
	}
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := server.ListenAndServeTLS(
			"", "",
		); err != nil && !errors.Is(
			err, http.ErrServerClosed,
		) {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// Listen for the interrupt signal.
	<-ctx.Done()
	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	timeToHandle := 5
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(timeToHandle)*time.Second)
	defer cancel()
	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %s", err)
	}
	// Close database connections
	if err := config.MongoDbPool.Client().Disconnect(ctx); err != nil {
		log.Printf("Error shutting down mongo connection: %v", err)
	}
	// Close redis connections
	if err := config.RedisPool.Close(); err != nil {
		log.Printf("Error shutting down redis connection: %v", err)
	}
	// notify user of shutdown
	log.Println("Server exiting")
}

func registerRoutes(engine *gin.Engine, sgCtx context.Context) {
	router := engine
	// no route handler
	router.NoRoute(func(ctx *gin.Context) {
		ctx.String(http.StatusNotFound,
			http.StatusText(http.StatusNotFound))
	})
	// main route handler
	router.GET(common.EmptyPath, func(ctx *gin.Context) {
		ctx.String(http.StatusOK, fmt.Sprintf("%s %s",
			config.Instance.AppName, config.Instance.AppVersion))
	})
	// health check handler
	healthConfig := healthcheckconfig.DefaultConfig()
	healthConfig.HealthPath = "/health"
	mongoStatus := healthcheck.NewMongoCheck(10,
		config.MongoDbPool.Client())
	redisStatus := healthcheck.NewRedisCheck(config.RedisPool)
	serverStatus := healthcheck.NewContextCheck(sgCtx, "signals")
	networkStatus := healthcheck.NewPingCheck("https://www.google.com",
		"GET", 1000, nil, nil)
	_ = ginhealthcheck.New(router, healthConfig, []healthcheck.Check{
		mongoStatus, &redisStatus, serverStatus, networkStatus,
	})
	// graphql playground and query handler
	graphQLRouterGroup := router.Group("graphql")
	graphQLRouterGroup.GET("/", func(ctx *gin.Context) {
		playground.Handler("GraphQL Playground", "/graphql/query").
			ServeHTTP(ctx.Writer, ctx.Request)
	})
	graphQLRouterGroup.POST("/query", func(ctx *gin.Context) {
		handler.NewDefaultServer(graph.NewExecutableSchema(
			graph.Config{Resolvers: &graph.Resolver{}})).
			ServeHTTP(ctx.Writer, ctx.Request)
	})
}
