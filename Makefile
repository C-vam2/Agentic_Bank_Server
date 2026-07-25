postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:18-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root agentic_bank_server

dropdb:
	docker exec -it postgres dropdb --username=root agentic_bank_server

migrateup:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/agentic_bank_server?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/agentic_bank_server?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/Agentic_Bank_Server/db/sqlc Store


.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock