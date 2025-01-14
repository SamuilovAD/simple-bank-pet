init:
	docker run --name postgres17 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine
createdb:
	docker exec -it postgres17 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres17 dropdb simple_bank
migrate-create:
	@echo "migrate create -ext sql -dir db/migration -seq $(name)"
migrate-up:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" up
migrate-up-next:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" up 1
migrate-down:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" down
migrate-down-last:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" down 1

sqlc:
	sqlc generate
testRunExample:
	go test github.com/SamuilovAD/simple-bank-pet/db/sqlc -run ^TestMain$
unit-tests-with-coverage:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/SamuilovAD/simple-bank-pet/db/sqlc Store
.PHONY: createdb dropdb migrate-create migrate-up migrate-up-next migrate-down migrate-down-last sqlc testRunExample server mock