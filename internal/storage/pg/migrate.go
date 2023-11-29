package pg

import (
	"embed"
	_ "embed"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/labstack/gommon/log"
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
		return fmt.Errorf("error creaate migrate: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error migrate: %w", err)
	}

	return nil
}
