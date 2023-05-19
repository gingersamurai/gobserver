.PHONY: test migrate local_postgres_init

run:
	go run gobserver/cmd/cli_app

build:
	go build -o ./build/gobserver gobserver

test:
	go test gobserver/...

migrate:
	goose -dir ./migrations up

local_postgres_init:
	docker run --rm --name local-postgres -p 5432:5432 -e POSTGRES_PASSWORD=12345678 -d postgres
