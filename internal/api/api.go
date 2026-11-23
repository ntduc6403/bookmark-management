package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/ntduc6403/bookmark-management/docs"
	"github.com/ntduc6403/bookmark-management/internal/handler"
	"github.com/ntduc6403/bookmark-management/internal/repository"
	"github.com/ntduc6403/bookmark-management/internal/service"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Engine is an interface that defines the methods for starting the server and serving HTTP requests.
type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

// engine is a struct that implements the Engine interface and holds the Gin engine and configuration.
type engine struct {
	app         *gin.Engine
	cfg         *Config
	redisClient *redis.Client
}

// NewEngine creates a new instance of the engine with the provided configuration.
func NewEngine(cfg *Config, redisClient *redis.Client) Engine {

	app := &engine{
		app:         gin.Default(),
		cfg:         cfg,
		redisClient: redisClient,
	}
	app.InitRoutes()
	return app
}

// Start starts the Gin server on the specified port from the configuration.
func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%s", e.cfg.AppPort))
}

// ServeHTTP serves HTTP requests using the Gin engine.
func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.app.ServeHTTP(w, req)
}

// InitRoutes initializes the routes for the Gin engine, including the health check and Swagger documentation endpoints.
func (e *engine) InitRoutes() {
	// Heath_CHECK
	redisHealthCheckRepo := repository.NewHealthCheck(e.redisClient)
	healthCheckSvc := service.NewHealthCheck(
		e.cfg.ServiceName,
		e.cfg.InstanceID,
		redisHealthCheckRepo,
	)
	healthCheckHandler := handler.NewHealthCheck(healthCheckSvc)

	e.app.GET("/health-check", healthCheckHandler.HealthCheck)

	// Shorten URL
	urlRepo := repository.NewURLStorage(e.redisClient)
	shortenURLSvc := service.NewShortenUrl(urlRepo)
	shortenURLHandler := handler.NewShortenURLHandler(shortenURLSvc)

	v1 := e.app.Group("/v1")
	{
		links := v1.Group("/links")
		{
			links.POST("/shorten", shortenURLHandler.Shorten)
		}
	}

	// Swagger
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}
