package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/merrickfox/go-scaffold/models"
	"net/http"
	"time"
)

type RefreshClaims struct {
	ExpiresAt int64  `json:"exp,omitempty"`
	Subject   string `json:"sub,omitempty"`
}

type AccessClaims struct {
	ExpiresAt   int64  `json:"exp,omitempty"`
	Subject     string `json:"sub,omitempty"`
	GivenName   string `json:"given_name,omitempty"`
	FamilyName  string `json:"family_name,omitempty"`
	Username    string `json:"username,omitempty"`
	UserRegDate int64 `json:"user_reg_date,omitempty"`
	Email       string `json:"email,omitempty"`
}

func (rc RefreshClaims) Valid() error {
	vErr := new(jwt.ValidationError)
	now := time.Now().Unix()
	fmt.Println("Expires at: ", rc.ExpiresAt)
	fmt.Println("now: ", now)
	if now > rc.ExpiresAt {
		delta := time.Unix(now, 0).Sub(time.Unix(rc.ExpiresAt, 0))
		vErr.Inner = fmt.Errorf("token is expired by %v", delta)
		vErr.Errors |= jwt.ValidationErrorExpired
	}

	if vErr.Errors == 0 {
		return nil
	}

	return vErr
}

func (rc RefreshClaims) GetExpiresAt() int64 {
	return rc.ExpiresAt
}

func (rc AccessClaims) Valid() error {
	return nil
}

func GenerateJwtPair(user models.UserDb, accessSigningKey, refreshSigningKey string) (*models.LoginResponse, *models.ServiceError) {
	uac := AccessClaims{
		ExpiresAt:   time.Now().Add(time.Minute * 15).Unix(),
		Subject:     user.Id,
		GivenName:   user.GivenName,
		FamilyName:  user.FamilyName,
		Username:    user.Username,
		UserRegDate: user.CreatedAt.Unix(),
		Email:       user.Email,
	}

	at, err := GenerateJwt([]byte(accessSigningKey), uac)
	if err != nil {
		return nil, err
	}

	urc := RefreshClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		Subject:   user.Id,
	}
	rt, err := GenerateJwt([]byte(refreshSigningKey), urc)
	if err != nil {
		return nil, err
	}

	tp := models.LoginResponse{
		AccessToken:  at,
		RefreshToken: rt,
	}

	return &tp, nil
}

func GenerateJwt(signingKey []byte, userClaims jwt.Claims) (string, *models.ServiceError) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = userClaims

	tokenString, err := token.SignedString(signingKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		se := models.NewServiceError(models.ServiceErrorInternalError, err.Error(), http.StatusInternalServerError, &err)
		return "", se
	}

	return tokenString, nil
}

func VerifyJwt(signingKey []byte, t string, claims jwt.Claims) (*jwt.Token, *models.ServiceError) {
	pt, err := jwt.ParseWithClaims(t, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return signingKey, nil
	})
	if err != nil {
		var se *models.ServiceError

		if jwtErr, ok := err.(*jwt.ValidationError); ok {
			if jwtErr.Errors == 16 {
				se = models.NewServiceError(models.ServiceErrorUnauthorised, err.Error(), http.StatusUnauthorized, &err)
			}
		} else if err == jwt.ErrSignatureInvalid {
			se = models.NewServiceError(models.ServiceErrorUnauthorised, err.Error(), http.StatusUnauthorized, &err)
		} else {
			se = models.NewServiceError(models.ServiceErrorInternalError, err.Error(), http.StatusInternalServerError, &err)
		}
		return nil, se
	}

	return pt, nil
}
