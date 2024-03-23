package postgress_pgxpool

/*
	Работа через pgx
*/

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"go_db/config"
	"go_db/internal/models"
	queries "go_db/internal/storage/postgress"
)

type PostgressStore struct {
	log  *slog.Logger
	conn *pgx.Conn
}

func GetStore(ctx context.Context, logger *slog.Logger) (*PostgressStore, error) {
	logger.Info("connecting DB postgress...")
	conn, err := pgx.Connect(ctx, config.GetConnectionString())
	//pool, err := pgxpool.New(ctx, config.GetConnectionString()) v5 не работает с pgxscan
	if err != nil {
		return nil, err
	}
	i := PostgressStore{log: logger, conn: conn}
	i.log.Info("DB postgress connected")
	return &i, nil
}

func (i PostgressStore) GetConnection(ctx context.Context) (*sql.DB, error) {
	//	i.pool.Get
	//	conn, err := i.pool.Acquire(ctx)
	return nil, errors.New("pgx: GetConnection not supported")
}

func (i PostgressStore) Release(ctx context.Context) {
	i.conn.Close(ctx)
	i.log.Info("DB postgress disconnected")
}

func (i PostgressStore) EmployeeCreate(ctx context.Context, empl models.Employee) (uint32, error) {
	i.log.Debug("pgx:create ", slog.String("Name", empl.Name))
	var id uint32
	//i.conn.BeginTx()
	row := i.conn.QueryRow(ctx, queries.QueryCreate, empl.Name)
	err := row.Scan(&id)
	return id, err
}

func (i PostgressStore) EmployeeGet(ctx context.Context, id uint32) (*models.Employee, error) {
	i.log.Debug("pgx:get ", slog.Any("ID", id))
	var empl models.Employee
	row := i.conn.QueryRow(ctx, queries.QueryGet, id)
	err := row.Scan(&empl.Id, &empl.Name)
	return &empl, err
}
func (i PostgressStore) EmployeeUpdate(ctx context.Context, empl models.Employee) error {
	i.log.Debug("pgx:update ", slog.Any("ID", empl.Id), slog.String("Name", empl.Name))
	commandTag, err := i.conn.Exec(ctx, queries.QueryUpdate, empl.Name, empl.Id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("not found")
	}
	return err
}
func (i PostgressStore) EmployeeDelete(ctx context.Context, id uint32) error {
	i.log.Debug("pgx:delete ", slog.Any("ID", id))
	commandTag, err := i.conn.Exec(ctx, queries.QueryDelete, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("not found")
	}
	return nil
}
func (i PostgressStore) EmployeeList(ctx context.Context) ([]*models.Employee, error) {
	i.log.Debug("pgx:list")
	var result []*models.Employee

	rows, err := i.conn.Query(ctx, queries.QueryList)
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
