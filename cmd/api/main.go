package main

import (
	"log"

	"github.com/ntduc6403/bookmark-management/internal/api"
	"github.com/ntduc6403/bookmark-management/pkg/redis"
)

// @title Bookmark API
// @version 1.0
// @description API documentation for the Bookmark  service.
// @host  localhost:8080
// @BasePath /
func main() {
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}
	client, err := redis.NewClient("")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	
	app := api.NewEngine(cfg, client)
	err = app.Start()
	if err != nil {
		panic(err)
	}

}
