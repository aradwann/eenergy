# Build stage
FROM golang:1.21.9-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.19 
WORKDIR /app
COPY --from=builder /app/main .
COPY .env .
COPY db/migrations ./db/migrations
COPY dev-certs ./dev-certs

EXPOSE 8080

ENTRYPOINT [ "/app/main" ]