# Project Name

## Overview

This project is designed to handle users, user accounts, and money transfers between them. It leverages the following technologies:
- **Golang**: For writing the main application and workers.
- **Postgres**: As the primary database, utilizing stored functions.
- **Redis**: For caching and as a job queue.
- **gRPC**: For public API interfacing, respecting versioning.
- **OpenTelemetry**: For collecting logs, traces, and metrics.

## Architecture

### Layers

1. **API Layer**: Handles gRPC requests and routes them to appropriate service handlers.
2. **Service Layer**: Contains business logic and application services.
3. **Data Layer**: Manages interaction with Postgres and Redis.
4. **Worker Layer**: Background job processing using Redis queues.
5. **Telemetry Layer**: Uses OpenTelemetry for logs, traces, and metrics.

### Components

#### API Layer (gRPC)

- **gRPC Server**: Handles incoming requests, performs authentication, and routes to service handlers.
- **Versioning**: Implements versioning in gRPC services for backward compatibility.

#### Service Layer

- **AccountService**: Manages user accounts and profile updates.
- **TransferService**: Handles money transfer logic, ensuring atomicity and consistency.
- **CacheService**: Interfaces with Redis for caching.
- **JobService**: Manages job scheduling and processing through Redis queues.

#### Data Layer

- **Postgres Repository**:
  - Interacts with the database using `database/sql` or an ORM like `GORM`.
  - Executes stored functions for complex transactions and business logic.
- **Redis Repository**:
  - Caching layer to reduce database load.
  - Job queue management using Redis.

#### Worker Layer

- **Worker**: Background Golang processes that:
  - Process queued jobs (e.g., sending notifications, processing transfers).
  - Use Redis for job queuing and state management.

#### Telemetry Layer

- **OpenTelemetry**:
  - Instrumentation for tracing, metrics, and logging.
  - Integrated with gRPC, database queries, and workers.

## Project Structure

```plaintext
.
├── api
│   └── grpc
│       ├── v1
│       │   ├── handlers
│       │   │   ├── account_handler.go
│       │   │   └── transfer_handler.go
│       │   └── proto
│       │       ├── account.proto
│       │       └── transfer.proto
│       ├── v2
│       │   ├── handlers
│       │   │   ├── account_handler.go
│       │   │   └── transfer_handler.go
│       │   └── proto
│       │       ├── account.proto
│       │       └── transfer.proto
├── cmd
│   ├── migrate
│   │   └── main.go
│   └── server
│       └── main.go
├── migrations
│   ├── 0001_create_accounts_table.up.sql
│   ├── 0001_create_accounts_table.down.sql
│   ├── 0002_create_transfers_table.up.sql
│   ├── 0002_create_transfers_table.down.sql
│   └── ...
├── service
│   ├── v1
│   │   ├── account_service.go
│   │   └── transfer_service.go
│   ├── v2
│   │   ├── account_service.go
│   │   └── transfer_service.go
│   ├── cache_service.go
│   └── job_service.go
├── repository
│   ├── postgres
│   │   ├── account_repo.go
│   │   └── transfer_repo.go
│   └── redis
│       ├── cache_repo.go
│       └── job_repo.go
├── worker
│   └── transfer_worker.go
├── telemetry
│   ├── logging.go
│   ├── metrics.go
│   └── tracing.go
└── Makefile
