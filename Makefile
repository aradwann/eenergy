
DB_SOURCE = "postgresql://root:secret@localhost:5432/eenergy?sslmode=disable"
MIGRATIONS_PATH = db/migrations
PROCS_PATH = db/procs

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root eenergy

dropdb:
	docker exec -it postgres15 dropdb eenergy

migrateup:
	go run db/scripts/migrate.go

migrateprocsup:
	migrate -path $(PROCS_PATH) -database $(DB_SOURCE) -verbose up

migrateup1:
	migrate -path $(MIGRATIONS_PATH) -database $(DB_SOURCE) -verbose up 1

migratedown:
	migrate -path $(MIGRATIONS_PATH) -database $(DB_SOURCE) -verbose down

migratedown1:
	migrate -path $(MIGRATIONS_PATH) -database $(DB_SOURCE) -verbose down 1

createmigration:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq "$(filter-out $@,$(MAKECMDGOALS))"

test: 
	go test -v -cover ./...

testci:
	go test -race -covermode atomic -coverprofile=covprofile ./...

server:
	go run main.go


protoc: 
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb  --grpc-gateway_opt paths=source_relative \
	--openapiv2_out=doc/swagger \
	--openapiv2_opt=allow_merge=true,merge_file_name=eenergy\
    proto/*.proto
	statik -src=./doc/swagger -dest=./doc -f

evans:
	evans --host localhost --port 9090 -r repl

.PHONEY: createdb dropdb migrateup migrateup1 migratedown migratedown1  test server  protoc evans migrateprocsup

