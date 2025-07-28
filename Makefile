DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable
postgres:
	docker run --name postgres17 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine
createdb:
	docker exec -it postgres17 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres17 dropdb simple_bank
migrate-create:
	@echo "migrate create -ext sql -dir db/migration -seq $(name)"
migrate-up:
	migrate -path db/migration -database "$(DB_URL)" up
migrate-up-next:
	migrate -path db/migration -database "$(DB_URL)" up 1
migrate-down:
	migrate -path db/migration -database "$(DB_URL)" down
migrate-down-last:
	migrate -path db/migration -database "$(DB_URL)" down 1
db_docs:
	dbdocs build doc/db.dbml --password "secret"
db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml
sqlc:
	sqlc generate
testRunExample:
	go test github.com/SamuilovAD/simple-bank-pet/db/sqlc -run ^TestMain$
test-with-coverage:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/SamuilovAD/simple-bank-pet/db/sqlc Store
	mockgen -package mockwk -destination worker/mock/distributor.go github.com/SamuilovAD/simple-bank-pet/worker TaskDistributor
proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto \
      --go_out=pb --go_opt=paths=source_relative \
      --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
      --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
      --openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank\
      proto/*.proto
evans:
	 evans \
       --host=localhost \
       --port 9090 \
       --proto service_simple_bank.proto \
     --path proto

redis:
	docker run --name redis -p 6379:6379 -d redis:8-alpine
lint:
	golangci-lint run
.PHONY: createdb dropdb migrate-create migrate-up migrate-up-next migrate-down migrate-down-last sqlc testRunExample test-with-coverage server mock db_docs db_schema proto redis lint