package redis

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"encoding/json"
	"log/slog"

	"github.com/go-redis/redis"

	"go_db/config"
	"go_db/internal/models"
	"go_db/internal/storage"
)

type RedisStore struct {
	cl  *redis.Client
	log *slog.Logger
	db  storage.Store
}

func GetStore(ctx context.Context, logger *slog.Logger, dbStore storage.Store) (*RedisStore, error) {
	rs := RedisStore{log: logger, db: dbStore}
	rs.cl = redis.NewClient(&redis.Options{
		Addr:     config.GetRedisAddr(),
		DB:       config.RedisEmployeeDb,
		Password: config.RedisPassword})
	rs.log.Debug("Redist client ok")
	return &rs, nil
}

func (i RedisStore) GetConnection(ctx context.Context) (*sql.DB, error) {
	return nil, errors.New("redis:GetConnection not supported")
}

func (i RedisStore) Release(ctx context.Context) {
	err := i.cl.Close()
	if err != nil {
		i.log.Error("redigo close error", err)
	}
	i.log.Debug("cache release")
}

func (i RedisStore) Set(ctx context.Context, empl models.Employee) error {
	data, err := json.Marshal(empl)
	if err != nil {
		return err
	}
	status := i.cl.Set(fmt.Sprint(empl.Id), data, config.RedisExpiration)
	return status.Err()
}

func (i RedisStore) Get(ctx context.Context, emplPtr *models.Employee) (bool, error) {
	sId := fmt.Sprint(emplPtr.Id)
	sCmd := i.cl.Get(sId)
	err := sCmd.Err()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	bytes, err := sCmd.Bytes()
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(bytes, emplPtr)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (i RedisStore) EmployeeCreate(ctx context.Context, empl models.Employee) (uint32, error) {
	i.log.Debug("redis:create ", slog.String("Name", empl.Name))
	var err error
	empl.Id, err = i.db.EmployeeCreate(ctx, empl)
	if err != nil {
		return 0, err
	}
	err = i.Set(ctx, empl)
	if err != nil {
		i.log.Error("Cache", slog.String("message", err.Error()))
	}

	return empl.Id, nil
}

func (i RedisStore) EmployeeGet(ctx context.Context, id uint32) (*models.Employee, error) {
	i.log.Debug("redis:get ", slog.Any("ID", id))
	empl := models.Employee{Id: id}
	found, err := i.Get(ctx, &empl)
	if found {
		i.log.Debug("cache found")
		return &empl, nil
	} else if err != nil {
		i.log.Error("cache ", slog.String("message", err.Error()))
	} else {
		i.log.Debug("cache not found")
	}
	return i.db.EmployeeGet(ctx, id)
}

func (i RedisStore) EmployeeUpdate(ctx context.Context, empl models.Employee) error {
	i.log.Debug("redis:update ", slog.Any("ID", empl.Id), slog.String("Name", empl.Name))
	err := i.db.EmployeeUpdate(ctx, empl)
	if err != nil {
		return err
	}
	err = i.Set(ctx, empl)
	if err != nil {
		i.log.Error("cache ", slog.String("message", err.Error()))
	}
	return nil
}

func (i RedisStore) EmployeeDelete(ctx context.Context, id uint32) error {
	i.log.Debug("redis:delete ", slog.Any("ID", id))

	err := i.db.EmployeeDelete(ctx, id)
	intCmd := i.cl.Del(fmt.Sprint(id))
	if intCmd.Err() != nil {
		i.log.Error("Cache delete", slog.String("message", intCmd.Err().Error()))
	}
	return err
}

func (i RedisStore) EmployeeList(ctx context.Context) ([]*models.Employee, error) {
	i.log.Debug("redis:list")
	return i.db.EmployeeList(ctx)
}
