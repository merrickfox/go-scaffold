package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/merrickfox/go-scaffold/config"
	"github.com/merrickfox/go-scaffold/resource"
	"net/http"
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
	e.POST("/login", h.login)
	e.POST("/refresh", h.refresh)

	r := e.Group("")
	r.Use(middleware.JWT([]byte(cfg.JwtAccessSecret)))
	r.POST("/reset-password", h.resetPassword)
	r.POST("/thing", someHandler)

}

func someHandler(c echo.Context) error {
	return c.String(http.StatusCreated, "yolo")
}