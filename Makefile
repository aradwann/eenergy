
DB_SOURCE = "postgresql://root:secret@localhost:5432/eenergy?sslmode=disable"
MIGRATIONS_PATH = db/migrations

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root eenergy

dropdb:
	docker exec -it postgres15 dropdb eenergy

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/aradwann/eenergy/db/store Store
	mockgen -package mockmail -destination mail/mock/sender.go github.com/aradwann/eenergy/mail EmailSender
	mockgen -package mockwk -destination worker/mock/distributor.go github.com/aradwann/eenergy/worker TaskDistributor
	mockgen -package mockwk -destination worker/mock/processor.go github.com/aradwann/eenergy/worker TaskProcessor

migrateup:
	go run db/scripts/migrate.go

migrateup1:
	migrate -path $(MIGRATIONS_PATH) -database $(DB_SOURCE) -verbose up 1

migratedown:
	migrate -path $(MIGRATIONS_PATH) -database $(DB_SOURCE) -verbose down

migratedown1:
	migrate -path $(MIGRATIONS_PATH) -database $(DB_SOURCE) -verbose down 1

createmigration:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq "$(filter-out $@,$(MAKECMDGOALS))"

test: 
	go test -short -v -cover ./...

testci:
	go test -short -race -covermode atomic -coverprofile=covprofile $(go list ./... | grep -v ./pb/)

server:
	go run main.go

protoc: 
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=eenergy \
	proto/*.proto

evans:
	evans --host localhost --port 9090 -r repl

.PHONEY: createdb dropdb migrateup migrateup1 migratedown migratedown1 test server protoc evans migrateprocsup

