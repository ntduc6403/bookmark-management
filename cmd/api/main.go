package main

import "github.com/ntduc6403/bookmark-management/internal/api"

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

	app := api.NewEngine(cfg)
	err = app.Start()
	if err != nil {
		panic(err)
	}
}
