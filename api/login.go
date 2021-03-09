package api

import (
	"github.com/labstack/echo/v4"
	"github.com/merrickfox/go-scaffold/crypto"
	"github.com/merrickfox/go-scaffold/jwt"
	"github.com/merrickfox/go-scaffold/models"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func (h *handler) login(c echo.Context) error {
	lr := new(models.LoginRequest)
	if err := c.Bind(lr); err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"email": lr.Email,
		"ip": c.RealIP(),
	}).Info("login request received")
	user, err := h.resource.FetchUserByEmail(lr.Email)
	if err != nil {
		return err.ToResponse(c)
	}

	if ok := crypto.CheckPasswordHash(lr.Password, user.HashedPassword); !ok {
		log.WithFields(log.Fields{
			"email": lr.Email,
			"ip": c.RealIP(),
		}).Info("failed login request")
		err = models.NewServiceError(models.ServiceErrorUnauthorised, "Incorrect user or password", http.StatusUnauthorized, nil)
		return err.ToResponse(c)
	}

	resp, err := jwt.GenerateJwtPair(*user, h.config.JwtAccessSecret, h.config.JwtRefreshSecret)
	if err != nil {
		return err.ToResponse(c)
	}

	at := new(http.Cookie)
	at.Name = "accessToken"
	at.Value = resp.AccessToken
	at.Expires = time.Now().Add(24 * time.Hour)
	at.SameSite = 3
	//at.Secure = true
	c.SetCookie(at)

	rt := new(http.Cookie)
	rt.Name = "refreshToken"
	rt.Value = resp.RefreshToken
	rt.Expires = time.Now().Add(24 * time.Hour)
	rt.SameSite = 3
	//rt.Secure = true
	c.SetCookie(rt)

	log.WithFields(log.Fields{
		"email": lr.Email,
		"ip": c.RealIP(),
	}).Info("successful login")
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