package api

import (
	"github.com/labstack/echo/v4"
	"github.com/merrickfox/go-scaffold/models"
	"net/http"
)

func (h *handler) register(c echo.Context) error {
	user := new(models.UserRequest)
	if err := c.Bind(user); err != nil {
		return err
	}

	err := user.Validate()
	if err != nil {
		return err.ToResponse(c)
	}

	dbu, err := user.ToDbStruct()
	if err != nil {
		return err.ToResponse(c)
	}

	err = h.resource.InsertUser(dbu)
	if err != nil {
		return err.ToResponse(c)
	}

	return c.String(http.StatusCreated, "")
}