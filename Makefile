postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=testPassword -d postgres:9.6.21-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:testPassword@localhost:5432/simple_bank?sslmode=disable" --verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:testPassword@localhost:5432/simple_bank?sslmode=disable" --verbose down

sqlc:
	sqlc generate

startdb:
	docker start postgres

stopdb:
	docker stop postgres

test:
	go test -v -cover ./...

.PHONY: createdb postgres dropdb migrateup migratedown sqlc startdb stopdb test