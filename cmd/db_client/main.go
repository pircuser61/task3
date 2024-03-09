package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"go_db/cmd/migrations"
	"go_db/config"
	"go_db/internal/models"

	//dbPackage "go_db/internal/storage/postgress/sql"
	//dbPackage "go_db/internal/storage/postgress/pgxpool"
	//dbPackage "go_db/internal/storage/postgress/pgx"
	//dbPackage "go_db/internal/storage/postgress/sqlx"
	//dbPackage "go_db/internal/storage/postgress/gorm"
	dbPackage "go_db/internal/storage/mongo"
)

func main() {
	var logLevel slog.Level
	switch config.LogLevel {
	case "Debug":
		logLevel = slog.LevelDebug
	case "Info":
		logLevel = slog.LevelInfo
	default:
		logLevel = slog.LevelError
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbInstanse, err := dbPackage.GetStore(ctx, logger)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	conn, err := dbInstanse.GetConnection(ctx)
	if err == nil {
		err = migrations.MakeMigrations(conn)
	}
	if err != nil {
		logger.Error(err.Error())
	}

	defer dbInstanse.Release(ctx)

	/*-------------------------------------------*/

	printList := func() {
		uList, err := dbInstanse.EmployeeList(ctx)
		if err != nil {
			fmt.Println("list error", err.Error())
		} else {
			for _, ui := range uList {
				fmt.Println(ui.Id, ui.Name)
			}
		}
	}

	ux := models.Employee{Name: "Вася"}

	ux.Id, err = dbInstanse.EmployeeCreate(ctx, ux)
	if err != nil {
		fmt.Println("Create error", err.Error())

	} else {
		fmt.Println("Created UserId", ux.Id)
	}

	printList()

	ux.Name = "Коля"
	err = dbInstanse.EmployeeUpdate(ctx, ux)
	if err != nil {
		fmt.Println("update error:", err.Error())
	}

	ptrEmpl, err := dbInstanse.EmployeeGet(ctx, ux.Id)
	if err != nil {
		fmt.Println("get error:", err.Error())
	} else {
		fmt.Println("get:", ptrEmpl.Id, ptrEmpl.Name)
	}

	dbInstanse.EmployeeDelete(ctx, ux.Id)
	if err != nil {
		fmt.Println("delete error:", err.Error())
	}
	printList()

}
