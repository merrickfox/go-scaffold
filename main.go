package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/merrickfox/go-scaffold/api"
	"github.com/merrickfox/go-scaffold/config"
	"github.com/merrickfox/go-scaffold/resource"
	log "github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("API started...")

	cfg := config.GetConfig()
	repo, closeRepo, err := resource.NewPostgresRepo(cfg)
	if err != nil {
		log.Fatal("could not start db")
	}
	defer closeRepo()

	if err := repo.Migrate(); err != nil {
		log.Fatalf("could not migrate : %v", err)
	}

	e := echo.New()
	//e.Use(middleware.Recover())
	api.Init(e, repo, cfg)

	// Start server
	e.Logger.Fatal(e.Start(":6969"))
}