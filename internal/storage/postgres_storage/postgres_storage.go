package postgres_storage

import (
	"context"
	"github.com/jackc/pgx/v5"
	"gobserver/internal/entity"
)

type PostgresStorage struct {
	conn *pgx.Conn
}

func NewPostgresStorage(connStr string) (*PostgresStorage, error) {
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{conn: conn}, nil
}

func (ps *PostgresStorage) AddCommandExecution(ctx context.Context, commandExecution entity.CommandExecution) (int64, error) {
	requestSQL := "INSERT INTO command_execution(project_path, command, exit_code, timestamp) VALUES ($1, $2, $3, $4) RETURNING ID"
	result := ps.conn.QueryRow(ctx, requestSQL,
		commandExecution.ProjectPath,
		commandExecution.Command,
		commandExecution.ExitCode,
		commandExecution.Timestamp)
	var id int64
	if err := result.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (ps *PostgresStorage) AddFileChange(ctx context.Context, fileChange entity.FileChange) (int64, error) {
	requestSQL := "INSERT INTO file_change(project_path, file_path, operation, timestamp) VALUES ($1, $2, $3, $4) RETURNING ID"
	result := ps.conn.QueryRow(ctx, requestSQL,
		fileChange.ProjectPath,
		fileChange.FilePath,
		fileChange.Operation,
		fileChange.Timestamp)
	var id int64
	if err := result.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
