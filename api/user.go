package api

import (
	"github.com/labstack/echo/v4"
	"github.com/merrickfox/go-scaffold/jwt"
	"net/http"
)

func (h *handler) user(c echo.Context) error {
	t := c.Get("user")
	u, se := jwt.GetUserFromToken(t)
	if se != nil {
		return se.ToResponse(c)
	}
	return c.JSON(http.StatusCreated, u)
}