package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/merrickfox/go-scaffold/crypto"
	"github.com/merrickfox/go-scaffold/jwt"
	"github.com/merrickfox/go-scaffold/models"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func (h *handler) resetPassword(c echo.Context) error {
	rp := new(models.ResetPasswordRequest)
	if err := c.Bind(rp); err != nil {
		return err
	}

	err := rp.Validate()
	if err != nil {
		return err.ToResponse(c)
	}


	fmt.Println(c.Get("user"))
	t := c.Get("user")

	cl, err := jwt.GetRawClaimsFromToken(t)
	if err != nil {
		return err.ToResponse(c)
	}

	email := cl["email"].(string)
	log.WithFields(log.Fields{
		"email": email,
		"ip": c.RealIP(),
	}).Info("password reset request")

	user, se := h.resource.FetchUserByEmail(email)
	if se != nil {
		return se.ToResponse(c)
	}

	if ok := crypto.CheckPasswordHash(rp.OldPassword, user.HashedPassword); !ok {
		log.WithFields(log.Fields{
			"email": email,
			"ip": c.RealIP(),
		}).Info("failed password reset request")
		err = models.NewServiceError(models.ServiceErrorUnauthorised, "Incorrect user or password", http.StatusUnauthorized, nil)
		return err.ToResponse(c)
	}

	hash, err2 := crypto.HashPassword(rp.NewPassword)
	if err2 != nil {
		se2 := models.NewServiceError(models.ServiceErrorInternalError, err2.Error(), http.StatusInternalServerError, &err2)
		return se2.ToResponse(c)
	}

	user.HashedPassword = hash
	n := time.Now()
	user.PasswordLastUpdated = &n

	se3 := h.resource.UpdatePassword(user)
	if se3 != nil {
		return se3.ToResponse(c)
	}
	return c.String(http.StatusNoContent, "")
}