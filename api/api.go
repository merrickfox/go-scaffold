package api

import (
	"github.com/labstack/echo/v4"
	"github.com/merrickfox/go-scaffold/config"
	"github.com/merrickfox/go-scaffold/resource"
)

type handler struct {
	resource        resource.Postgres
	config config.Config
}

func Init(e *echo.Echo, repo resource.Postgres, cfg config.Config) {
	h := &handler{
		resource: repo,
		config: cfg,
	}

	e.POST("/register", h.register)
	//e.POST("/login", h.login)
}
