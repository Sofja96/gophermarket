package pg

import (
	_ "embed"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/labstack/gommon/log"
	_ "github.com/labstack/gommon/log"
	"os"
)

////go:embed migrations/*.sql
//var migrateFS embed.FS

func migrateDatabase(dsn string) error {
	//d, err := iofs.New(migrateFS, "migrations")
	//if err != nil {
	//	return fmt.Errorf("error on creating iofs: %w", err)
	//}
	//
	//m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	//if err != nil {
	//	return fmt.Errorf("error on creating migration: %w", err)
	//}
	//
	//if err := m.Up(); err != nil && err != migrate.ErrNoChange {
	//	return fmt.Errorf("error on making migration: %w", err)
	//}

	//return nil

	file := "file://./internal/storage/pg/migrations/"
	init, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("init script is not found: %w", err)
	}
	m, err := migrate.New(
		string(init), dsn)
	if err != nil {
		return fmt.Errorf("unable to create a migrator: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("unable to migrate: %w", err)

	}

	log.Infof("Migration done. Current schema version: %v\n")
	return nil
}
