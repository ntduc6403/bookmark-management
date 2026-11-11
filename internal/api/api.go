package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/ntduc6403/bookmark-management/docs"
	"github.com/ntduc6403/bookmark-management/internal/handler"
	"github.com/ntduc6403/bookmark-management/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

type engine struct {
	app *gin.Engine
	cfg *Config
}

func NewEngine(cfg *Config) Engine {

	app := &engine{
		app: gin.Default(),
		cfg: cfg,
	}
	app.InitRoutes()
	return app
}

func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%s", e.cfg.AppPort))
}
func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.app.ServeHTTP(w, req)
}
func (e *engine) InitRoutes() {
	healthCheckSvc := service.NewHealthCheck(
		e.cfg.ServiceName,
		e.cfg.InstanceID,
	)
	healthCheckHandler := handler.NewHealthCheck(healthCheckSvc)

	e.app.GET("/health-check", healthCheckHandler.HealthCheck)
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
