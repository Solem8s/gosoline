package db

import (
	"fmt"
	"github.com/applike/gosoline/pkg/mon"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type MigrationSettings struct {
	Application    string `cfg:"application" default:"{app_name}"`
	Path           string `cfg:"path"`
	PrefixedTables bool   `cfg:"prefixed_tables" default:"false"`
	Enabled        bool   `cfg:"enabled" default:"false"`
}

func runMigrations(logger mon.Logger, settings Settings, db *sqlx.DB) {
	if !settings.Migrations.Enabled || settings.Migrations.Path == "" {
		return
	}

	driverFactory, err := GetDriverFactory(settings.Driver)

	if err != nil {
		logger.Panicf(err, "could not get driver factory for %s", settings.Driver)
	}

	migrationsTable := "schema_migrations"

	if settings.Migrations.PrefixedTables {
		application := strings.ToLower(settings.Migrations.Application)
		application = strings.Replace(application, "-", "_", -1)
		migrationsTable = fmt.Sprintf("%s_schema_migrations", application)
	}

	driver, err := driverFactory.GetMigrationDriver(db.DB, settings.Uri.Database, migrationsTable)

	if err != nil {
		logger.Panic(err, "could not get migration driver")
	}

	ex, _ := os.Executable()
	migrationsAbsoluteFilePath := filepath.Join(filepath.Dir(ex), settings.Migrations.Path)

	filePath := fmt.Sprintf("file://%s", migrationsAbsoluteFilePath)

	if runtime.GOOS == "windows" {
		filePath = migrationsAbsoluteFilePath
	}

	m, err := migrate.NewWithDatabaseInstance(filePath, settings.Driver, driver)

	if err != nil {
		logger.Panic(err, "could not initialize migrator for db migrations")
	}

	err = m.Up()

	if err == migrate.ErrNoChange {
		logger.Info("no db migrations to run")
		return
	}

	if err != nil {
		logger.Panic(err, "could not run db migrations")
	}
}
