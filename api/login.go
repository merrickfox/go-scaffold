package api

import (
	"github.com/labstack/echo/v4"
	"github.com/merrickfox/go-scaffold/crypto"
	"github.com/merrickfox/go-scaffold/jwt"
	"github.com/merrickfox/go-scaffold/models"
	"net/http"
)

func (h *handler) login(c echo.Context) error {
	lr := new(models.LoginRequest)
	if err := c.Bind(lr); err != nil {
		return err
	}

	user, err := h.resource.FetchUserByEmail(lr.Email)
	if err != nil {
		return err.ToResponse(c)
	}

	if ok := crypto.CheckPasswordHash(lr.Password, user.HashedPassword); !ok {
		err = models.NewServiceError(models.ServiceErrorUnauthorised, "Incorrect user or password", http.StatusUnauthorized, nil)
		err.ToResponse(c)
	}

	resp, err := jwt.GenerateJwtPair(*user, h.config.JwtAccessSecret, h.config.JwtRefreshSecret)
	if err != nil {
		return err.ToResponse(c)
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *handler) refresh(c echo.Context) error {
	rr := new(models.RefreshRequest)
	if err := c.Bind(rr); err != nil {
		return err
	}

	rc := jwt.RefreshClaims{}
	pt, err := jwt.VerifyJwt([]byte(h.config.JwtRefreshSecret), rr.RefreshToken, &rc)
	if err != nil {
		return err.ToResponse(c)
	}

	ptClaims, ok := pt.Claims.(*jwt.RefreshClaims)
	if !ok {
		err = models.NewServiceError(models.ServiceErrorInternalError, "internal error", http.StatusInternalServerError, nil)
		return err.ToResponse(c)
	}

	user, err := h.resource.FetchUserById(ptClaims.Subject)
	if err != nil {
		return err.ToResponse(c)
	}

	resp, err := jwt.GenerateJwtPair(*user, h.config.JwtAccessSecret, h.config.JwtRefreshSecret)

	return c.JSON(http.StatusOK, resp)
}