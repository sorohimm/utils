package migrate

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/pkg/errors"
)

func NewMigrator(cfg *Config) (*Migrator, error) {
	return &Migrator{
		url:    cfg.Postgres.URL,
		schema: cfg.Postgres.Schema,
	}, nil
}

type Migrator struct {
	url    string
	schema string
}

func (o *Migrator) Up(url string) error {
	db, err := sql.Open("postgres", o.url)
	if err != nil {
		return fmt.Errorf("open db error: %w", err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("driver instance error: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(url, "postgres", driver)
	if err != nil {
		return fmt.Errorf("migrate instance error: %w", err)
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}

		return fmt.Errorf("up error: %w", err)
	}

	return nil
}
