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
	i := PostgressStore{log: logger, db: db}
	i.log.Info("DB postgress connected")
	return &i, nil
}

func (i PostgressStore) GetConnection(_ context.Context) (*sql.DB, error) {
	return i.db, nil
}

func (i PostgressStore) Release(_ context.Context) {
	i.db.Close()
	i.log.Info("DB postgress disconnected")
}

func (i PostgressStore) EmployeeCreate(ctx context.Context, empl models.Employee) (uint32, error) {
	var id uint32
	i.log.Debug("sql:create ", slog.String("Name", empl.Name))
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
	tx, err := i.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, queries.QueryCreate, empl.Name)
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	tx.Commit()

	return id, nil
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
		return nil, err
	}
	defer rows.Close()

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
