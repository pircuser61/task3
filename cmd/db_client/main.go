package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	migrate_goose "go_db/cmd/migrations/goose"
	migrate_migrate "go_db/cmd/migrations/migrate"
	"go_db/config"
	"go_db/internal/models"
	"go_db/internal/storage"

	//dbPackage "go_db/internal/storage/postgress/sql"
	//dbPackage "go_db/internal/storage/postgress/pgxpool"
	//dbPackage "go_db/internal/storage/postgress/pgx"
	//dbPackage "go_db/internal/storage/postgress/sqlx"
	//dbPackage "go_db/internal/storage/postgress/gorm"
	//dbPackage "go_db/internal/storage/postgress/go-pg"
	dbPackage "go_db/internal/storage/mongo"

	//cachePackage "go_db/internal/storage/go-redis"
	cachePackage "go_db/internal/storage/redigo"
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

	var dbInstanse, cacheInstance storage.Store
	var err error
	dbInstanse, err = dbPackage.GetStore(ctx, logger)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer dbInstanse.Release(ctx)

	conn, err := dbInstanse.GetConnection(ctx)
	if err == nil {
		logger.Debug("===== Migrations =====")
		err = migrate_goose.MakeMigrations(conn)
		if err != nil {
			logger.Error(err.Error())
		} else {
			logger.Debug("Goose: ok")
		}
		err = migrate_migrate.MakeMigrations(conn, logger)
		if err != nil {
			logger.Error(err.Error())
		} else {
			logger.Debug("Migrate: ok")
		}

	} else if err != nil {
		logger.Error(err.Error())
		logger.Debug("===== Migrations skipped =====")
	}

	cacheInstance, err = cachePackage.GetStore(ctx, logger, dbInstanse)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer cacheInstance.Release(ctx)

	dbInstanse = cacheInstance
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

	err = dbInstanse.EmployeeDelete(ctx, ux.Id)
	if err != nil {
		fmt.Println("delete error:", err.Error())
	}
	printList()

}
