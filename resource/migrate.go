package resource

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	post "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
)

func (p *Postgres) Migrate() error {
	driver, err := post.WithInstance(p.Db.DB, &post.Config{})
	if err != nil {
		return fmt.Errorf("could not create database driver : %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://resource/migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("error instantiating migrate : %w", err)
	}

	version, dirty, _ := m.Version()
	log.Printf("Database version %d, dirty %t", version, dirty)

	log.Printf("Starting migration")

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("an error occurred while syncing the database.. %w", err)
	}

	log.Println("Migration complete")
	return nil
}