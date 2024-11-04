init:
	docker run --name postgres17 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine
createdb:
	docker exec -it postgres17 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres17 dropdb simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" up
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" down
sqlc:
	sqlc generate
testRunExample:
	go test github.com/SamuilovAD/simple-bank-pet/db/sqlc -run ^TestMain$
unit-tests-with-coverage:
	go test -v -cover ./...
.PHONY: createdb dropdb migrateup migratedown sqlc testRunExample