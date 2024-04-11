package postgress_sqlx

/*
	Работа через database/sql
*/

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	//_ "github.com/lib/pq"

	"go_db/config"
	"go_db/internal/models"
	queries "go_db/internal/storage/postgress"
)

type PostgressStore struct {
	log *slog.Logger
	db  *sqlx.DB
}

func GetStore(ctx context.Context, logger *slog.Logger) (*PostgressStore, error) {
	logger.Info("connecting DB postgress...")
	/*
	   	db, err := sqlx.Connect("pgx",  config.GetConnectionString())
	       if err != nil {
	          return nil, err
	       }
	*/
	db, err := sqlx.Open("pgx", config.GetConnectionString())
	if err != nil {
		return nil, err
	}
	i := PostgressStore{log: logger, db: db}
	i.log.Info("DB postgress connected")
	return &i, nil
}

func (i PostgressStore) GetConnection(_ context.Context) (*sql.DB, error) {
	return i.db.DB, nil
}

func (i PostgressStore) Release(_ context.Context) {
	err := i.db.Close()
	if err != nil {
		i.log.Error("sqlx close error", err)
	}
	i.log.Info("DB postgress disconnected")
}

func (i PostgressStore) EmployeeCreate(_ context.Context, empl models.Employee) (uint32, error) {
	var id uint32
	i.log.Debug("sqlx:create ", slog.String("Name", empl.Name))
	/*
		sqlResult, err := i.db.Exec("INSERT INTO employee (Name) VALUES ($1)",
				empl.Name)
					// "github.com/lib/pq" : LastInsertId is not supported by this driver
		id64, err = sqlResult.LastInsertId()
		if err != nil {
			return 0, err
		}
			id := uint32(id64)
	*/
	//i.db.BeginTx()
	row := i.db.QueryRow(queries.QueryCreate, empl.Name)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (i PostgressStore) EmployeeGet(_ context.Context, id uint32) (*models.Employee, error) {
	i.log.Debug("sqlx:get ", slog.Any("ID", id))
	var empl models.Employee
	row := i.db.QueryRow(queries.QueryGet, id)
	err := row.Scan(&empl.Id, &empl.Name)
	return &empl, err
}
func (i PostgressStore) EmployeeUpdate(_ context.Context, empl models.Employee) error {
	i.log.Debug("sqlx:update ", slog.Any("ID", empl.Id), slog.String("Name", empl.Name))
	sqlResult, err := i.db.Exec(queries.QueryUpdate,
		empl.Name, empl.Id)
	if err != nil {
		return err
	}
	rowCount, err := sqlResult.RowsAffected()
	if err != nil {
		return err
	}
	if rowCount != 1 {
		return errors.New("not found")
	}
	return err
}
func (i PostgressStore) EmployeeDelete(_ context.Context, id uint32) error {
	i.log.Debug("sqlx:delete ", slog.Any("ID", id))
	sqlResult, err := i.db.Exec(queries.QueryDelete, id)
	if err != nil {
		return err
	}
	rowCount, err := sqlResult.RowsAffected()
	if err != nil {
		return err
	}
	if rowCount != 1 {
		return errors.New("not found")
	}
	return nil
}
func (i PostgressStore) EmployeeList(_ context.Context) ([]*models.Employee, error) {
	i.log.Debug("sqlx:list")
	var result []*models.Employee

	err := i.db.Select(&result, queries.QueryList)
	return result, err
}
