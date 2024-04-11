package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	migrate_goose "go_db/cmd/migrations/goose"
	"go_db/config"
	"go_db/internal/models"
	"go_db/internal/storage"
	dbPackage "go_db/internal/storage/postgress/sql"
	cachePackage "go_db/internal/storage/redigo"
)

func main() {
	//fmt.Printf("bad format for vet 2", 12)
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
			return
		} else {
			logger.Debug("Goose: ok")
		}
	}
	cacheInstance, err = cachePackage.GetStore(ctx, logger, dbInstanse)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer cacheInstance.Release(ctx)

	dbInstanse = cacheInstance

	id, err := dbInstanse.EmployeeCreate(ctx, models.Employee{Name: "TheFirst"})
	if err != nil {
		logger.Error(err.Error())
		return
	}
	logger.Debug("Created", slog.Int("Id", int(id)))

	list := func(rw http.ResponseWriter, _ *http.Request) {
		logger.Debug("get list")
		resp, err := dbInstanse.EmployeeList(context.Background())
		makeResp(rw, 0, resp, err)
	}

	/*
	   func get(rw http.ResponseWriter, _ *http.Request) {
	   	fmt.Fprintf(rw, "HI THERE")
	   }

	   func update(rw http.ResponseWriter, _ *http.Request) {
	   	fmt.Fprintf(rw, "HI THERE")
	   }

	   func delete(rw http.ResponseWriter, _ *http.Request) {
	   	fmt.Fprintf(rw, "HI THERE")
	   }
	*/

	http.HandleFunc("/", list)
	err = http.ListenAndServe(config.AppAddr, nil)
	if err != nil {
		panic(err)
	}
}

func makeResp(rw http.ResponseWriter, tm int64, body any, err error) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	rw.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	type response struct {
		Error  bool
		ErrMsg string
		Body   any
	}

	if err != nil {
		_ = json.NewEncoder(rw).Encode(response{Error: true, ErrMsg: err.Error()})
	} else {
		_ = json.NewEncoder(rw).Encode(response{Body: body})
	}
}
