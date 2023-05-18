-- +goose Up
-- +goose StatementBegin
CREATE TABLE command_execution(
    id SERIAL PRIMARY KEY,
    project_path TEXT,
    command TEXT,
    exit_code INTEGER,
    timestamp TIMESTAMP
);

CREATE TABLE file_change(
    id SERIAL PRIMARY KEY,
    project_path TEXT,
    file_path TEXT,
    operation TEXT,
    timestamp TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE command_execution;
DROP TABLE file_change;
-- +goose StatementEnd
