postgres:
	docker run --name postgres14 -e POSTGRES_PASSWORD=123456 -e POSTGRES_USER=root -p 5432:5432 -d postgres:14-alpine

createdb:
	docker exec -it postgres14 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres14 dropdb simple_bank

migrateup:
	migrate --path db/migration --database="postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate --path db/migration --database="postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	docker run --rm -v "%cd%:/src" -w /src kjconroy/sqlc init
	docker run --rm -v $(pwd):/src -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres, createdb, dropdb, migrateup, migratedown, sqlc, test
