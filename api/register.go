package api

import (
	"github.com/labstack/echo/v4"
	"github.com/merrickfox/go-scaffold/models"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (h *handler) register(c echo.Context) error {
	user := new(models.UserRequest)
	if err := c.Bind(user); err != nil {
		return err
	}


	err := user.Validate()
	if err != nil {
		log.WithFields(log.Fields{
			"email": user.Email,
			"username": user.Username,
			"given_name": user.GivenName,
			"family_name": user.FamilyName,
			"ip": c.RealIP(),
		}).Info("failed registration")
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

	log.WithFields(log.Fields{
		"email": user.Email,
		"ip": c.RealIP(),
	}).Info("new registration")
	return c.String(http.StatusCreated, "")
}