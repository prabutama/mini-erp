IMAGE_REGISTRY ?= ghcr.io
IMAGE_OWNER ?= your-ghcr-owner
IMAGE_TAG ?= dev

.PHONY: help build test tidy clean docker-build docker-push postgres-up postgres-down postgres-logs migrate-identity-up migrate-organization-up migrate-operations-up migrate-resource-up migrate-reporting-up run-identity run-organization run-operations run-resource run-reporting run-api-gateway

help:
	@echo "make build  - build all Go services"
	@echo "make test   - test all Go services"
	@echo "make tidy   - tidy all Go modules"
	@echo "make clean  - remove local build output directory"
	@echo "make docker-build - build all service images"
	@echo "make docker-push  - push all service images"
	@echo "make postgres-up   - start local PostgreSQL"
	@echo "make postgres-down - stop local PostgreSQL"
	@echo "make run-identity     - run Identity gRPC service"
	@echo "make run-organization - run Organization gRPC service"
	@echo "make run-operations   - run Operations gRPC service"
	@echo "make run-resource     - run Resource gRPC service"
	@echo "make run-reporting    - run Reporting gRPC service"
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

docker-build:
	docker build -f deploy/docker/go-service.Dockerfile --build-arg SERVICE=api-gateway -t $(IMAGE_REGISTRY)/$(IMAGE_OWNER)/mini-erp-api-gateway:$(IMAGE_TAG) .
	docker build -f deploy/docker/go-service.Dockerfile --build-arg SERVICE=identity -t $(IMAGE_REGISTRY)/$(IMAGE_OWNER)/mini-erp-identity:$(IMAGE_TAG) .
	docker build -f deploy/docker/go-service.Dockerfile --build-arg SERVICE=organization -t $(IMAGE_REGISTRY)/$(IMAGE_OWNER)/mini-erp-organization:$(IMAGE_TAG) .
	docker build -f deploy/docker/go-service.Dockerfile --build-arg SERVICE=operations -t $(IMAGE_REGISTRY)/$(IMAGE_OWNER)/mini-erp-operations:$(IMAGE_TAG) .
	docker build -f deploy/docker/go-service.Dockerfile --build-arg SERVICE=resource -t $(IMAGE_REGISTRY)/$(IMAGE_OWNER)/mini-erp-resource:$(IMAGE_TAG) .
	docker build -f deploy/docker/go-service.Dockerfile --build-arg SERVICE=reporting -t $(IMAGE_REGISTRY)/$(IMAGE_OWNER)/mini-erp-reporting:$(IMAGE_TAG) .

docker-push:
	docker push $(IMAGE_REGISTRY)/$(IMAGE_OWNER)/mini-erp-api-gateway:$(IMAGE_TAG)
	docker push $(IMAGE_REGISTRY)/$(IMAGE_OWNER)/mini-erp-identity:$(IMAGE_TAG)
	docker push $(IMAGE_REGISTRY)/$(IMAGE_OWNER)/mini-erp-organization:$(IMAGE_TAG)
	docker push $(IMAGE_REGISTRY)/$(IMAGE_OWNER)/mini-erp-operations:$(IMAGE_TAG)
	docker push $(IMAGE_REGISTRY)/$(IMAGE_OWNER)/mini-erp-resource:$(IMAGE_TAG)
	docker push $(IMAGE_REGISTRY)/$(IMAGE_OWNER)/mini-erp-reporting:$(IMAGE_TAG)

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

migrate-operations-up:
	migrate -path services/operations/migrations -database "postgres://mini_erp:mini_erp@localhost:5432/operations_db?sslmode=disable" up

migrate-resource-up:
	migrate -path services/resource/migrations -database "postgres://mini_erp:mini_erp@localhost:5432/resource_db?sslmode=disable" up

migrate-reporting-up:
	migrate -path services/reporting/migrations -database "postgres://mini_erp:mini_erp@localhost:5432/reporting_db?sslmode=disable" up

run-identity:
	cd services/identity && IDENTITY_DATABASE_URL="postgres://mini_erp:mini_erp@localhost:5432/identity_db?sslmode=disable" go run ./cmd/identity

run-organization:
	cd services/organization && ORGANIZATION_DATABASE_URL="postgres://mini_erp:mini_erp@localhost:5432/organization_db?sslmode=disable" go run ./cmd/organization

run-operations:
	cd services/operations && OPERATIONS_DATABASE_URL="postgres://mini_erp:mini_erp@localhost:5432/operations_db?sslmode=disable" go run ./cmd/operations

run-resource:
	cd services/resource && RESOURCE_DATABASE_URL="postgres://mini_erp:mini_erp@localhost:5432/resource_db?sslmode=disable" go run ./cmd/resource

run-reporting:
	cd services/reporting && REPORTING_DATABASE_URL="postgres://mini_erp:mini_erp@localhost:5432/reporting_db?sslmode=disable" go run ./cmd/reporting

run-api-gateway:
	cd services/api-gateway && IDENTITY_GRPC_ADDR="localhost:50051" ORGANIZATION_GRPC_ADDR="localhost:50052" OPERATIONS_GRPC_ADDR="localhost:50053" RESOURCE_GRPC_ADDR="localhost:50054" REPORTING_GRPC_ADDR="localhost:50055" go run ./cmd/api-gateway
