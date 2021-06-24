package resource

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

type Postgres struct {
	Db *sqlx.DB
}

func NewPostgresRepo(cfg config.Config) (Postgres, func() error, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s "+
		"dbname=%s sslmode=disable",
		cfg.DbHost,
		cfg.DbPort,
		cfg.DbUser,
		cfg.DbPassword,
		cfg.DbName)
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("could not open db: ", err)
	}

	localPg := Postgres{Db: db}

	return localPg, localPg.Db.Close, err
}
