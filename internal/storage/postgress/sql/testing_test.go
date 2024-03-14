package postgress_sql

import (
	"context"
	"log/slog"
	"os"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

type PostgressMockStore struct {
	store PostgressStore
	mock  sqlmock.Sqlmock
}

func setUp(t *testing.T) (PostgressMockStore, error) {
	var ms PostgressMockStore
	var err error
	ms.store.log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	ms.store.db, ms.mock, err = sqlmock.New()
	return ms, err
}

func (i *PostgressMockStore) tearDown() {
	i.store.Release(context.Background())
}
