package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/merrickfox/go-scaffold/models"
	"net/http"
	"time"
)

func GenerateJwt(signingKey []byte, userClaims map[string]interface{}) (string, *models.ServiceError) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = userClaims["authorized"]
	claims["client"] = userClaims["client"]
	claims["exp"] = time.Now().Add(time.Hour * 24 * 365).Unix()

	tokenString, err := token.SignedString(signingKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		se := models.NewServiceError(models.ServiceErrorInternalError, err.Error(), http.StatusInternalServerError, &err)
		return "", se
	}

	return tokenString, nil
}

func VerifyJwt(signingKey []byte, t string) *models.ServiceError {
	_, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return signingKey, nil
	})
	if err != nil {
		var se *models.ServiceError
		if err == jwt.ErrSignatureInvalid {
			se = models.NewServiceError(models.ServiceErrorUnauthorised, err.Error(), http.StatusUnauthorized, &err)
		} else {
			se = models.NewServiceError(models.ServiceErrorInternalError, err.Error(), http.StatusInternalServerError, &err)

		}
		fmt.Println("Something Went Wrong: ", err.Error())
		return se
	}

	return nil
}