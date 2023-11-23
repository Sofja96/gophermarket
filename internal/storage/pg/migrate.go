package pg

import (
	"embed"
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
		return fmt.Errorf("error on creating migration: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error on making migration: %w", err)
	}

	return nil

	//file := "file://./internal/storage/pg/migrations/"
	//init, err := os.ReadFile(file)
	//if err != nil {
	//	return fmt.Errorf("init script is not found: %w", err)
	//}
	//m, err := migrate.New(
	//	"file://./internal/storage/pg/migrations/", dsn)
	//if err != nil {
	//	return fmt.Errorf("unable to create a migrator: %w", err)
	//}
	//if err := m.Up(); err != nil && err != migrate.ErrNoChange {
	//	return fmt.Errorf("unable to migrate: %w", err)
	//
	//}
	//migrator, err := migrate.NewMigrator(conn, "schema_version")
	//if err != nil {
	//	log.Fatalf("Unable to create a migrator: %v\n", err)
	//}
	//
	//err = migrator.LoadMigrations("./migrations")
	//if err != nil {
	//	log.Fatalf("Unable to load migrations: %v\n", err)
	//}
	//
	//err = migrator.Migrate()
	//if err != nil {
	//	log.Fatalf("Unable to migrate: %v\n", err)
	//}
	//
	//ver, err := migrator.GetCurrentVersion()
	//if err != nil {
	//	log.Fatalf("Unable to get current schema version: %v\n", err)
	//}
	//
	//log.Infof("Migration done. Current schema version: %v\n", ver)
	return nil
}
