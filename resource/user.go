package resource

import (
	"fmt"
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

func (p Postgres) FetchUserByEmail(email string) (*models.UserDb, *models.ServiceError) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select().
		Columns("*").
		From("users").
		Where(sq.Eq{"email": email})
	var user models.UserDb
	sql, args, _ := query.ToSql()
	err := p.Db.Get(&user, sql, args...)
	if err != nil {
		fmt.Println(err)
		return nil, models.NewServiceError(models.ServiceErrorUnauthorised, "Incorrect user or password", http.StatusUnauthorized, &err)
	}

	return &user, nil
}

func (p Postgres) FetchUserById(id string) (*models.UserDb, *models.ServiceError) {
	fmt.Println(id)
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select().
		Columns("*").
		From("users").
		Where(sq.Eq{"id": id})
	var user models.UserDb
	sql, args, _ := query.ToSql()
	err := p.Db.Get(&user, sql, args...)
	if err != nil {
		fmt.Println(err)
		return nil, models.NewServiceError(models.ServiceErrorUnauthorised, "Incorrect user or password", http.StatusUnauthorized, &err)
	}

	return &user, nil
}