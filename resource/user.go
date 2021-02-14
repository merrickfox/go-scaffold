package resource

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/merrickfox/go-scaffold/models"
	"net/http"
)

func (p Postgres) InsertUser(user *models.UserDb) *models.ServiceError {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Insert("users").
		Columns("username", "email", "given_name", "family_name", "hashed_password").
		Values(user.Username, user.Email, user.GivenName, user.FamilyName, user.HashedPassword).
		Suffix("RETURNING \"id\"").
		RunWith(p.Db.DB)
	err := query.QueryRow().Scan(&user.Id)
	if err, ok := err.(*pq.Error); ok {
		if err.Code == "23505" {
			return models.NewServiceError(models.ServiceErrorUnprocessable, "user already exists", http.StatusUnprocessableEntity, nil)
		}
	}
	if err != nil {
		return models.NewServiceError(models.ServiceErrorInternalError, err.Error(), http.StatusInternalServerError, &err)
	}

	return nil
}