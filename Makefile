# Variables
DB_NAME=eenergy
DB_USER=root
DB_PASS=secret
DB_HOST=localhost
DB_PORT=5432
DB_SOURCE=postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable
MIGRATIONS_PATH=migrations
CERT_PATH=dev-certs/

# Setup
init:
# create dev certs path
	mkdir -p ${CERT_PATH} 

# Database Operations
# TODO: impl handy DB operations

# Migration
migrateup:
	go run cmd/migrate/main.go

migrateup1:
	migrate -path $(MIGRATIONS_PATH) -database $(DB_SOURCE) -verbose up 1

migratedown:
	migrate -path $(MIGRATIONS_PATH) -database $(DB_SOURCE) -verbose down

migratedown1:
	migrate -path $(MIGRATIONS_PATH) -database $(DB_SOURCE) -verbose down 1

createmigration:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq "$(filter-out $@,$(MAKECMDGOALS))"

# Mocks
mock: # TODO: edit to comply with the new architecture
	mockgen -package mockdb -destination repository/mock/store.go github.com/aradwann/eenergy/repository/store Store
	mockgen -package mockmail -destination mail/mock/sender.go github.com/aradwann/eenergy/mail EmailSender
	mockgen -package mockwk -destination worker/mock/distributor.go github.com/aradwann/eenergy/worker TaskDistributor
	mockgen -package mockwk -destination worker/mock/processor.go github.com/aradwann/eenergy/worker TaskProcessor

# Testing
test: 
	go test -short -v -cover ./...

testci:
	go test -short -race -covermode atomic -coverprofile=covprofile $$(go list ./... | grep -v /pb$$)

# Run Server
server:
	go run main.go

# Protocol Buffers
protoc: 
	protoc --proto_path=api/grpc/v1/proto \
  	--go_out=paths=source_relative:api/grpc/v1/handlers \
  	--go-grpc_out=paths=source_relative:api/grpc/v1/handlers \
  	--grpc-gateway_out=paths=source_relative:api/grpc/v1/handlers \
	--grpc-gateway_opt generate_unbound_methods=true \
	--openapiv2_out=assets/doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=eenergy,generate_unbound_methods=true \
  	api/grpc/v1/proto/**/*.proto 

#TODO: try to use grpc-gateway config through external config file 

# gRPC Client
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


.PHONY: createdb dropdb migrateup migrateup1 migratedown migratedown1 test server protoc evans gen-cert init

