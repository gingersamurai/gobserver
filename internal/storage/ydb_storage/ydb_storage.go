package ydb_storage

import (
	"context"
	"github.com/ydb-platform/ydb-go-sdk/v3"
)

type YDBStorage struct {
	db *ydb.Driver
}

func NewYDBStorage(connStr string) (*YDBStorage, error) {
	ctx := context.Background()
	db, err := ydb.Open(ctx, connStr)
	if err != nil {
		return nil, err
	}
	return &YDBStorage{db: db}, nil
}
