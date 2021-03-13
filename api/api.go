package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/merrickfox/go-scaffold/config"
	"github.com/merrickfox/go-scaffold/models"
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

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAccessControlAllowCredentials, echo.HeaderAuthorization, echo.HeaderAccessControlAllowMethods},
		AllowCredentials: true,
	}))
	e.Use(ErrorHandler)
	e.POST("/register", h.register)
	e.POST("/login", h.login, middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(1)))
	e.POST("/refresh", h.refresh)


	r := e.Group("")
	r.Use(middleware.JWT([]byte(cfg.JwtAccessSecret)))
	r.POST("/reset-password", h.resetPassword)
	r.POST("/thing", someHandler)
	r.GET("/user", h.user)

}

func someHandler(c echo.Context) error {
	return c.String(http.StatusCreated, "yolo")
}

func ErrorHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			ee, ok := err.(*echo.HTTPError)
			if ok && ee.Code == 400 && ee.Message == "missing or malformed jwt" {
				println("here")
				se := models.NewServiceError(models.ServiceErrorUnauthorised, "unauthorised", http.StatusUnauthorized, nil)
				se.ToResponse(c)
			} else {
				c.Error(err)
			}
		}
		return nil
	}
}
