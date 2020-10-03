package db

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4/database"
	_ "github.com/mattn/go-sqlite3"
	"net/url"
)

const DriverSqlite = "sqlite3"

func init() {
	connectionFactories[DriverSqlite] = NewSqliteDriverFactory()
}

func NewSqliteDriverFactory() DriverFactory {
	return &sqliteDriverFactory{}
}

type sqliteDriverFactory struct{}

func (m *sqliteDriverFactory) GetDSN(settings Settings) string {
	dsn := url.URL{
		User: url.UserPassword(settings.Uri.User, settings.Uri.Password),
		Host: fmt.Sprintf("file:%s", settings.Uri.Host),
		Path: settings.Uri.Database,
	}

	qry := dsn.Query()
	dsn.RawQuery = qry.Encode()

	uri := dsn.String()

	return uri[4:]
}

func (m *sqliteDriverFactory) GetMigrationDriver(db *sql.DB, database string, migrationsTable string) (database.Driver, error) {
	return nil, nil
}
