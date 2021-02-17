package models

import (
	"fmt"
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/merrickfox/go-scaffold/crypto"
	"net/http"
	"time"
)

type UserRequest struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	GivenName  string `json:"givenName"`
	FamilyName string `json:"familyName"`
	Password   string `json:"password"`
}

type UserDb struct {
	Id                  string     `db:"id"`
	Username            string     `db:"username"`
	Email               string     `db:"email"`
	GivenName           string     `db:"given_name"`
	FamilyName          string     `db:"family_name"`
	HashedPassword      string     `db:"hashed_password"`
	PasswordLastUpdated *time.Time `db:"password_last_updated"`
	CreatedAt           *time.Time `db:"created_at"`
	UpdatedAt           *time.Time `db:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshRequest struct {
	RefreshToken    string `json:"refreshToken"`
}

func (ur UserRequest) Validate() *ServiceError {
	fmt.Println("user: ")
	fmt.Println(ur)
	err := validation.ValidateStruct(&ur,
		validation.Field(&ur.Email, validation.Required, is.Email, validation.Length(1, 50)),
		validation.Field(&ur.Username, validation.Required, validation.Length(3, 20)),
		validation.Field(&ur.GivenName, validation.Required, validation.Length(1, 50)),
		validation.Field(&ur.FamilyName, validation.Required, validation.Length(1, 50)),
		validation.Field(&ur.Password, validation.Required, validation.Length(5, 80)),
	)

	if err != nil {
		se := NewServiceError(ServiceErrorInvalidRequestBody, err.Error(), http.StatusBadRequest, &err)
		return se
	}
	return nil
}

func (ur UserRequest) ToDbStruct() (*UserDb, *ServiceError) {
	udb := UserDb{
		Username:   ur.Username,
		Email:      ur.Email,
		GivenName:  ur.GivenName,
		FamilyName: ur.FamilyName,
	}

	hash, err := crypto.HashPassword(ur.Password)
	if err != nil {
		se := NewServiceError(ServiceErrorInternalError, err.Error(), http.StatusInternalServerError, &err)
		return nil, se
	}

	udb.HashedPassword = hash
	return &udb, nil
}
