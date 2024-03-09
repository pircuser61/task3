package postgress_pgxpool

/*
	Работа через pgxpool + pgxscan
*/

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	//"github.com/jackc/pgx/v5 похоже не совместим с pgxscan ..."
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"

	"go_db/config"
	"go_db/internal/models"
	queries "go_db/internal/storage/postgress"
)

type PostgressStore struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func GetStore(ctx context.Context, logger *slog.Logger) (*PostgressStore, error) {
	logger.Info("connecting DB postgress...")
	pool, err := pgxpool.Connect(ctx, config.GetConnectionString())
	//pool, err := pgxpool.New(ctx, config.GetConnectionString()) v5 не работает с pgxscan
	if err != nil {
		return nil, err
	}
	i := PostgressStore{log: logger, pool: pool}
	i.log.Info("DB postgress connected")
	return &i, nil
}

func (i PostgressStore) GetConnection(ctx context.Context) (*sql.DB, error) {
	//	i.pool.Get
	//	conn, err := i.pool.Acquire(ctx)
	return nil, errors.New("not supported")
}

func (i PostgressStore) Release(ctx context.Context) {
	i.pool.Close()
	i.log.Info("DB postgress disconnected")
}

func (i PostgressStore) EmployeeCreate(ctx context.Context, empl models.Employee) (uint32, error) {
	i.log.Debug("pgxpool:create ", slog.String("Name", empl.Name))
	var id uint32
	//tx, err:= i.pool.BeginTx()

	err := pgxscan.Get(ctx, i.pool, &id, queries.QueryCreate, empl.Name)
	return id, err
}

func (i PostgressStore) EmployeeGet(ctx context.Context, id uint32) (*models.Employee, error) {
	i.log.Debug("pgxpool:get ", slog.Any("ID", id))
	var empl models.Employee
	err := pgxscan.Get(ctx, i.pool, &empl, queries.QueryGet, id)
	return &empl, err
}
func (i PostgressStore) EmployeeUpdate(ctx context.Context, empl models.Employee) error {
	i.log.Debug("pgxpool:update ", slog.Any("ID", empl.Id), slog.String("Name", empl.Name))
	commandTag, err := i.pool.Exec(ctx, queries.QueryUpdate, empl.Name, empl.Id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("not found")
	}
	return err
}
func (i PostgressStore) EmployeeDelete(ctx context.Context, id uint32) error {
	i.log.Debug("pgxpool:delete ", slog.Any("ID", id))
	commandTag, err := i.pool.Exec(ctx, queries.QueryDelete, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("not found")
	}
	return nil
}
func (i PostgressStore) EmployeeList(ctx context.Context) ([]*models.Employee, error) {
	i.log.Debug("pgxpool:list")
	var result []*models.Employee
	if err := pgxscan.Select(ctx, i.pool, &result, queries.QueryList); err != nil {
		return nil, err
	}
	return result, nil
}
