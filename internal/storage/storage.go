package storage

import (
	"context"
	"database/sql"
	"go_db/internal/models"
)

type Store interface {
	GetConnection(context.Context) (*sql.DB, error)
	EmployeeCreate(context.Context, models.Employee) (uint32, error)
	EmployeeGet(context.Context, uint32) (*models.Employee, error)
	EmployeeUpdate(context.Context, models.Employee) error
	EmployeeDelete(context.Context, uint32) error
	EmployeeList(context.Context) ([]*models.Employee, error)
	Release(context.Context)
}
