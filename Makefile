.PHONY: help build test tidy clean postgres-up postgres-down postgres-logs migrate-identity-up migrate-organization-up run-identity run-organization run-api-gateway

help:
	@echo "make build  - build all Go services"
	@echo "make test   - test all Go services"
	@echo "make tidy   - tidy all Go modules"
	@echo "make clean  - remove local build output directory"
	@echo "make postgres-up   - start local PostgreSQL"
	@echo "make postgres-down - stop local PostgreSQL"
	@echo "make run-identity     - run Identity gRPC service"
	@echo "make run-organization - run Organization gRPC service"
	@echo "make run-api-gateway  - run API Gateway"

build:
	go build ./services/api-gateway/...
	go build ./services/identity/...
	go build ./services/organization/...
	go build ./services/operations/...
	go build ./services/resource/...
	go build ./services/reporting/...

test:
	go test ./services/api-gateway/...
	go test ./services/identity/...
	go test ./services/organization/...
	go test ./services/operations/...
	go test ./services/resource/...
	go test ./services/reporting/...

tidy:
	cd services/api-gateway && go mod tidy
	cd services/identity && go mod tidy
	cd services/organization && go mod tidy
	cd services/operations && go mod tidy
	cd services/resource && go mod tidy
	cd services/reporting && go mod tidy

clean:
	@if exist bin rmdir /s /q bin

postgres-up:
	docker compose up -d postgres

postgres-down:
	docker compose down

postgres-logs:
	docker compose logs -f postgres

migrate-identity-up:
	migrate -path services/identity/migrations -database "postgres://mini_erp:mini_erp@localhost:5432/identity_db?sslmode=disable" up

migrate-organization-up:
	migrate -path services/organization/migrations -database "postgres://mini_erp:mini_erp@localhost:5432/organization_db?sslmode=disable" up

run-identity:
	cd services/identity && IDENTITY_DATABASE_URL="postgres://mini_erp:mini_erp@localhost:5432/identity_db?sslmode=disable" go run ./cmd/identity

run-organization:
	cd services/organization && ORGANIZATION_DATABASE_URL="postgres://mini_erp:mini_erp@localhost:5432/organization_db?sslmode=disable" go run ./cmd/organization

run-api-gateway:
	cd services/api-gateway && IDENTITY_GRPC_ADDR="localhost:50051" ORGANIZATION_GRPC_ADDR="localhost:50052" go run ./cmd/api-gateway
