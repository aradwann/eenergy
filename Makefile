
DB_SOURCE = "postgresql://root:secret@localhost:5432/eenergy?sslmode=disable"
MIGRATIONS_PATH = db/migrations
CONFIG_PATH=dev-certs/

init:
	mkdir -p ${CONFIG_PATH}

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
	go test -short -race -covermode atomic -coverprofile=covprofile $$(go list ./... | grep -v /pb$$)

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
	evans --host localhost --port 9091 -r repl

# certs to be added to certs directory for local development purposes
gen-cert:
# Create the CA private key
	openssl genrsa -out ca-key.pem 2048
# Create a self-signed CA certificate
	openssl req -x509 -new -nodes -key ca-key.pem -days 3650 -out ca.pem -subj "/C=US/ST=NY/L=NYC/O=eenergy/CN=CA"
# Create the server private key
	openssl genrsa -out server.key 2048
# Create the server CSR
	openssl req -new -key server.key -out server.csr -subj "/C=US/ST=NY/L=NYC/O=eenergy/CN=server"
# Sign the server CSR with the CA certificate
	openssl x509 -req -in server.csr -CA ca.pem -CAkey ca-key.pem -CAcreateserial -out server.crt -days 365 -sha256
# Create the client private key
	openssl genrsa -out client.key 2048
# Create the client CSR
	openssl req -new -key client.key -out client.csr -subj "/C=US/ST=NY/L=NYC/O=eenergy/CN=client"
# Sign the client CSR with the CA certificate
	openssl x509 -req -in client.csr -CA ca.pem -CAkey ca-key.pem -CAcreateserial -out client.crt -days 365 -sha256

	mv *.pem *.csr *.crt *.srl *.key ${CONFIG_PATH}


.PHONEY: createdb dropdb migrateup migrateup1 migratedown migratedown1 test server protoc evans gen-cert init

