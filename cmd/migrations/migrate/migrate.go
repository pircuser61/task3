package migrate

import (
	"database/sql"
	"fmt"
	"log/slog"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type MyLogger struct {
	logger *slog.Logger
}

// Printf is like fmt.Printf
func (i MyLogger) Printf(format string, v ...interface{}) {
	i.logger.Debug(fmt.Sprintf(format, v...))

}

// Verbose should return true when verbose logging output is wanted
func (i MyLogger) Verbose() bool { return true }

func MakeMigrations(db *sql.DB, logger *slog.Logger) error {
	/*
		m, err := migrate.New(
			"github://mattes:personal-access-token@mattes/migrate_test",
			"postgres://localhost:5432/database?sslmode=enable")
	*/
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	dir := "D:/projects/go_db/migrations/migrate/"
	migrationsPath := fmt.Sprintf("file://%s", dir)
	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres", driver)
	if err != nil {
		return err
	}
	myLogger := MyLogger{logger}
	m.Log = myLogger
	logger.Debug("migrate:DOWN")
	err = m.Down()
	if err != nil {
		return err
	}
	logger.Debug("migrate:UP")
	err = m.Up()
	if err != nil {
		return err
	}
	return nil
}
