package postgress_gopg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-pg/pg/v10"

	"go_db/config"
	"go_db/internal/models"
)

/* CUSTOM TYPE EXAMPLE

import (
	"github.com/go-pg/pg/v10/types"
)
type Time struct {
	time.Time
}
func (tm Time) AppendValue(b []byte, flags int) ([]byte, error)
func (tm *Time) ScanValue(rd types.Reader, n int) error

var _ types.ValueScanner = (*Time)(nil)
var _ types.ValueAppender = (*Time)(nil)

type Event struct {
	Id   int
	Time Time `pg:"type:time"`
}

func main() {
	err := db.Model((*Event)(nil)).CreateTable(&orm.CreateTableOptions{
		Temp: true,
	})
    evt := new(Event)
	db.Model(&evt}).Insert()
	err = db.Model(evt).Select()
}
*/

type PostgressStore struct {
	log *slog.Logger
	db  *pg.DB
}

func GetStore(_ context.Context, logger *slog.Logger) (*PostgressStore, error) {
	/*
		opt, err := pg.ParseURL(config.GetConnectionString())
		db := pg.Connect(opt)
	*/
	logger.Info("connecting DB postgress...")
	pgOpt := config.GetConnectionOpt()
	opt := pg.Options{
		Addr:     fmt.Sprintf("%s:%s", pgOpt.Host, pgOpt.Port),
		Database: "empl",
		User:     pgOpt.User,
		Password: pgOpt.Passw,
	}
	db := pg.Connect(&opt)

	i := PostgressStore{log: logger, db: db}
	i.log.Info("go-pg: DB postgress connected")
	return &i, nil
}

func (i PostgressStore) GetConnection(_ context.Context) (*sql.DB, error) {
	return nil, errors.New("go-pg: Not supported")
}

func (i PostgressStore) Release(_ context.Context) {
	err := i.db.Close()
	if err != nil {
		i.log.Error("go-pg: close error", err)
	}
	i.log.Info("go-pg:Close")
}

func (i PostgressStore) EmployeeCreate(ctx context.Context, empl models.Employee) (newiD uint32, err error) {

	i.log.Debug("go-pg:create ", slog.String("Name", empl.Name))
	tx, err := i.db.Begin()
	if err != nil {
		return 0, err
	}
	//defer tx.Rollback()
	defer func() {
		closeErr := tx.Close() // Close calls Rollback if the tx has not already been committed or rolled back.
		if err == nil {
			err = closeErr
		}
	}()
	_, err = i.db.Model(&empl).Insert()
	if err != nil {
		return 0, err
	}

	if empl.Id > 0 {
		err := tx.Commit()
		return empl.Id, err
	}
	return 0, errors.New("empl_id == 0")
}

func (i PostgressStore) EmployeeGet(_ context.Context, id uint32) (*models.Employee, error) {
	i.log.Debug("go-pg:get ", slog.Any("ID", id))
	empl := models.Employee{Id: id}
	err := i.db.Model(&empl).WherePK().Select()
	return &empl, err
}

func (i PostgressStore) EmployeeUpdate(_ context.Context, empl models.Employee) error {
	i.log.Debug("go-pg:update ", slog.Any("ID", empl.Id), slog.String("Name", empl.Name))
	// db.Model(empl).Set("name = ?Name").Where("id = ?id").Update() Обновить одно поле
	//  db.Model(book).WherePK().Update() по PrimaryKey
	result, err := i.db.Model(&empl).
		Where("empl_id = $1", empl.Id).
		Update()
	if err != nil {
		return err
	}

	if result.RowsAffected() != 1 {
		return errors.New("not found")
	}
	return nil
}

func (i PostgressStore) EmployeeDelete(_ context.Context, id uint32) error {
	i.log.Debug("go-pg:delete ", slog.Any("ID", id))

	empl := models.Employee{Id: id}
	result, err := i.db.Model(&empl).
		WherePK().
		Delete()
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return errors.New("not found")
	}
	return errors.New("TEST ERROR")
}

func (i PostgressStore) EmployeeList(ctx context.Context) ([]*models.Employee, error) {
	i.log.Debug("go-pg:list")

	var emplList []*models.Employee

	err := i.db.Model(&emplList).Select()
	return emplList, err

	//return emplList, nil
}
