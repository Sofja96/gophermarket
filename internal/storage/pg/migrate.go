package pg

import (
	"embed"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrateFS embed.FS

func migrateDatabase(dsn string) error {
	d, err := iofs.New(migrateFS, "migrations")
	if err != nil {
		return fmt.Errorf("error on creating iofs: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		return fmt.Errorf("error on creating migration: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error on making migration: %w", err)
	}

	return nil

}
