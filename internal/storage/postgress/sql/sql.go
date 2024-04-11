package postgress_sql

/*
	Работа через database/sql
*/

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	//	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"

	"go_db/config"
	"go_db/internal/models"
	queries "go_db/internal/storage/postgress"
)

type PostgressStore struct {
	log *slog.Logger
	db  *sql.DB
}

func GetStore(ctx context.Context, logger *slog.Logger) (*PostgressStore, error) {
	logger.Info("connecting DB postgress...")
	db, err := sql.Open("postgres", config.GetConnectionString())
	if err != nil {
		return nil, err
	}
	logger.Info("ping DB...")
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	i := PostgressStore{log: logger, db: db}
	/*
		    DB.SetConnMaxIdleTime
			db.SetConnMaxLifetime(time.Minute * 3)
			db.SetMaxOpenConns(10)
			db.SetMaxIdleConns(10)
	*/

	i.log.Info("sql:DB postgress connected")
	return &i, nil
}

func GetMockStore(db *sql.DB, logger *slog.Logger) (*PostgressStore, error) {
	logger.Info("Mocking DB postgress...")
	i := PostgressStore{log: logger, db: db}
	return &i, nil
}

func (i PostgressStore) GetConnection(_ context.Context) (*sql.DB, error) {
	return i.db, nil
}

func (i PostgressStore) Release(_ context.Context) {
	err := i.db.Close()
	if err != nil {
		i.log.Error("sql close error", err)
	}
	i.log.Info("DB postgress disconnected")
}

func (i PostgressStore) EmployeeCreate(ctx context.Context, empl models.Employee) (uint32, error) {
	var id uint32
	i.log.Debug("sql:create ", slog.String("Name", empl.Name))
	/*
		sqlResult, err := i.db.Exec(queries.QueryCreate, empl.Name)
					// "github.com/lib/pq" : LastInsertId is not supported by this driver
		id64, err = sqlResult.LastInsertId()
		if err != nil {
			return 0, err
		}
			id := uint32(id64)
	*/
	tx, err := i.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		i.log.Error("sql", slog.String("create error", err.Error()))
		return 0, err
	}
	defer func() {
		err := tx.Rollback()
		if err != nil {
			// Что делать ?
		}
	}()

	row := tx.QueryRowContext(ctx, queries.QueryCreate, empl.Name)
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	err = tx.Commit()

	return id, err
}

func (i PostgressStore) EmployeeGet(_ context.Context, id uint32) (*models.Employee, error) {
	i.log.Debug("sql:get ", slog.Any("ID", id))
	var empl models.Employee
	row := i.db.QueryRow(queries.QueryGet, id)
	err := row.Scan(&empl.Id, &empl.Name)
	return &empl, err
}

func (i PostgressStore) EmployeeUpdate(_ context.Context, empl models.Employee) error {
	i.log.Debug("sql:update ", slog.Any("ID", empl.Id), slog.String("Name", empl.Name))
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
	i.log.Debug("sql:delete ", slog.Any("ID", id))
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

func (i PostgressStore) EmployeeList(ctx context.Context) ([]*models.Employee, error) {
	i.log.Debug("sql:list")
	var result []*models.Employee
	// rows, err := i.db.Query(ctx, queries.QueryList)
	rows, err := i.db.QueryContext(ctx, queries.QueryList)
	if err != nil {
		i.log.Debug("sql", slog.String("list ERROR", err.Error()))
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			i.log.Error("rows close error", err)
		}
	}()

	for rows.Next() {
		empl := models.Employee{}
		err := rows.Scan(&empl.Id, &empl.Name)
		if err != nil {
			return nil, err
		}
		result = append(result, &empl)
	}
	return result, nil
}
