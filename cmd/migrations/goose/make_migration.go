package migrations

import (
	"database/sql"

	//_ "github.com/lib/pq"
	goose "github.com/pressly/goose/v3"
)

func MakeMigrations(db *sql.DB) error {
	/* При запуске в докере нужен путь вида ./migrations/goose/ */
	/* я хз */
	return goose.Up(db, "./migrations/goose/")
}
