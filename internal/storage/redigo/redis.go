package redigo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"encoding/json"
	"log/slog"

	"github.com/gomodule/redigo/redis"

	"go_db/config"
	"go_db/internal/models"
	"go_db/internal/storage"
)

type RedisStore struct {
	conn redis.Conn
	log  *slog.Logger
	db   storage.Store
}

func GetStore(ctx context.Context, logger *slog.Logger, dbStore storage.Store) (*RedisStore, error) {
	rs := RedisStore{log: logger, db: dbStore}
	opt := redis.DialConnectTimeout(10 * time.Second)
	//conn := redis.DialTimeout("tcp", config.RedisAddr, 10 *time.Second)
	//depricated

	conn, err := redis.DialContext(ctx, "tcp", config.RedisAddr, opt)
	if err != nil {
		return nil, err
	}

	rs.conn = conn
	rs.log.Debug("Redist connect ok")
	return &rs, nil
}

func (i RedisStore) GetConnection(ctx context.Context) (*sql.DB, error) {
	return nil, errors.New("redis: GetConnection not supported")
}

func (i RedisStore) Release(ctx context.Context) {
	i.conn.Close()
	i.log.Debug("cache release")
}

func (i RedisStore) Set(ctx context.Context, empl models.Employee) error {
	data, err := json.Marshal(empl)
	if err != nil {
		return err
	}
	_, err = i.conn.Do("SET", fmt.Sprint(empl.Id), data)
	if err != nil {
		return fmt.Errorf("SET error: %w", err)
	}
	seconds := int(config.RedisExpiration.Seconds())
	_, err = i.conn.Do("EXPIRE", fmt.Sprint(empl.Id), seconds)
	if err != nil {
		return fmt.Errorf("SET EXPIRE error: %w", err)
	}
	return nil
}

func (i RedisStore) Get(ctx context.Context, emplPtr *models.Employee) (bool, error) {
	sId := fmt.Sprint(emplPtr.Id)
	result, err := redis.Bytes(i.conn.Do("GET", sId))
	if err == redis.ErrNil {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("redis GET error: %w", err)
	}
	err = json.Unmarshal(result, emplPtr)
	if err != nil {
		return false, fmt.Errorf("unmarshal error: %w", err)
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
		i.log.Error("Cache", slog.String("Set", err.Error()))
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
	_, cachErr := i.conn.Do("DEL", fmt.Sprint(id))
	if err != nil {
		i.log.Error("Cache delete", slog.String("message", cachErr.Error()))
	}
	return err
}

func (i RedisStore) EmployeeList(ctx context.Context) ([]*models.Employee, error) {
	i.log.Debug("redis:list")
	return i.db.EmployeeList(ctx)
}
