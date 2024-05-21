# Build stage
FROM golang:1.22.2-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/server/main.go

# Run stage
FROM alpine:3.19 
WORKDIR /app
COPY --from=builder /app/main .
COPY .env .
COPY migrations migrations
COPY dev-certs ./dev-certs

EXPOSE 8080

ENTRYPOINT [ "/app/main" ]