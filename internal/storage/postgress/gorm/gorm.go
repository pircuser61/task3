package postgres_gorm

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go_db/config"
	"go_db/internal/models"
)

type PostgressStore struct {
	log *slog.Logger
	db  *gorm.DB
}

func GetStore(_ context.Context, logger *slog.Logger) (*PostgressStore, error) {
	logger.Info("connecting DB postgress...")
	db, err := gorm.Open(postgres.Open(config.GetConnectionString()), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	i := PostgressStore{log: logger, db: db}
	i.log.Info("DB postgress connected")
	return &i, nil
}

func (i PostgressStore) GetConnection(_ context.Context) (*sql.DB, error) {
	return i.db.DB()
}

func (i PostgressStore) Release(_ context.Context) {
	i.log.Info("DB postgress disconnected")
}

func (i PostgressStore) EmployeeCreate(ctx context.Context, empl models.Employee) (uint32, error) {

	i.log.Debug("gorm:create ", slog.String("Name", empl.Name))
	tx := i.db.Begin()
	defer tx.Rollback()
	result := tx.Create(&empl)
	if result.Error != nil {
		return 0, result.Error
	}
	if empl.Id > 0 {
		tx.Commit()
		return empl.Id, nil
	}
	return 0, errors.New("empl_id == 0")
}

func (i PostgressStore) EmployeeGet(_ context.Context, id uint32) (*models.Employee, error) {
	i.log.Debug("gorm:get ", slog.Any("ID", id))
	var empl models.Employee
	result := i.db.First(&empl, id)

	return &empl, result.Error
}

func (i PostgressStore) EmployeeUpdate(_ context.Context, empl models.Employee) error {
	i.log.Debug("gorm:update ", slog.Any("ID", empl.Id), slog.String("Name", empl.Name))
	// i.db.Model(&empl).Update("Name", empl.Name) Обновить одно поле
	result := i.db.Model(&empl).Updates(empl)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return errors.New("not found")
	}
	return nil
}

func (i PostgressStore) EmployeeDelete(_ context.Context, id uint32) error {
	i.log.Debug("gorm:delete ", slog.Any("ID", id))
	result := i.db.Delete(&models.Employee{}, id)
	if result.RowsAffected != 1 {
		return errors.New("not found")
	}
	return result.Error
}

func (i PostgressStore) EmployeeList(ctx context.Context) ([]*models.Employee, error) {
	i.log.Debug("gorm:list")
	var emplList []*models.Employee
	result := i.db.Find(&emplList)
	return emplList, result.Error
}
