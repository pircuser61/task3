package postgress_sql

import (
	"context"
	"go_db/internal/models"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestCreate(t *testing.T) {

	t.Run("create", func(t *testing.T) {
		// arrange
		ms, err := setUp(t)
		defer ms.tearDown()
		if err != nil {
			t.Errorf("mock error")
		}
		empl := models.Employee{Id: 22, Name: "TestName"}

		ms.mock.ExpectBegin()
		expectQuery := regexp.QuoteMeta("INSERT INTO employee (Name) VALUES ($1) RETURNING empl_id")
		ms.mock.
			ExpectQuery(expectQuery).
			WithArgs(empl.Name).
			WillReturnRows(ms.mock.NewRows([]string{"id"}).AddRow(empl.Id))

		ms.mock.ExpectCommit()
		//ms.mock.ExpectRollback()

		// act
		id, err := ms.store.EmployeeCreate(context.Background(), empl)

		// assert
		if err != nil {
			t.Errorf("Unexpected error: %s ", err)
		}
		if id != empl.Id {
			t.Errorf("Id don't match")
		}
		if err := ms.mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
}

func TestGet(t *testing.T) {

	t.Run("get", func(t *testing.T) {
		// arrange
		ms, err := setUp(t)
		defer ms.tearDown()
		if err != nil {
			t.Errorf("mock error")
		}
		refEmpl := models.Employee{Id: 17, Name: "TestName"}
		expectQuery := regexp.QuoteMeta("SELECT * FROM employee WHERE empl_id = $1")
		ms.mock.
			ExpectQuery(expectQuery).
			WithArgs(refEmpl.Id).
			WillReturnRows(ms.mock.NewRows([]string{"id", "name"}).AddRow(refEmpl.Id, refEmpl.Name))

		// act
		resultEmpl, err := ms.store.EmployeeGet(context.Background(), refEmpl.Id)

		// assert
		if err != nil {
			t.Errorf("Unexpected error: %s ", err)
		}
		if *resultEmpl != refEmpl {
			t.Errorf("Empl has diff")
		}
		if err := ms.mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("update", func(t *testing.T) {
		// arrange
		ms, err := setUp(t)
		defer ms.tearDown()
		if err != nil {
			t.Errorf("mock error")
		}
		var empl models.Employee
		empl.Name = "TestName"

		expectQuery := regexp.QuoteMeta("UPDATE employee set name=$1 WHERE empl_id = $2")
		ms.mock.
			ExpectExec(expectQuery).
			WithArgs(empl.Name, empl.Id).
			WillReturnResult(sqlmock.NewResult(0, 1)) // no insert id, 1 affected row

		// act
		err = ms.store.EmployeeUpdate(context.Background(), empl)

		// assert
		if err != nil {
			t.Errorf("Unexpected error: %s ", err)
		}
		if err := ms.mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("delete", func(t *testing.T) {
		// arrange
		ms, err := setUp(t)
		defer ms.tearDown()
		if err != nil {
			t.Errorf("mock error")
		}
		id := uint32(22)

		expectQuery := regexp.QuoteMeta("DELETE FROM employee WHERE empl_id = $1")
		ms.mock.
			ExpectExec(expectQuery).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1)) // no insert id, 1 affected row

		// act
		err = ms.store.EmployeeDelete(context.Background(), id)

		// assert
		if err != nil {
			t.Errorf("Unexpected error: %s ", err)
		}
		if err := ms.mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
}
