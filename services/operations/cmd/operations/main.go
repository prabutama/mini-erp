package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	operationsevents "github.com/isapr/mini-erp/services/operations/internal/adapters/events"
	"github.com/isapr/mini-erp/services/operations/internal/adapters/grpcserver"
	operationsnats "github.com/isapr/mini-erp/services/operations/internal/adapters/nats"
	"github.com/isapr/mini-erp/services/operations/internal/adapters/postgres"
	"github.com/isapr/mini-erp/services/operations/internal/application"
	"github.com/isapr/mini-erp/services/operations/internal/ports"
	"google.golang.org/grpc"
)

func main() {
	databaseURL := env("OPERATIONS_DATABASE_URL", "postgres://mini_erp:mini_erp@localhost:5432/operations_db?sslmode=disable")
	addr := env("OPERATIONS_GRPC_ADDR", ":50053")

	pool, err := postgres.OpenPool(context.Background(), databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	services := postgres.NewServiceDefinitionRepository(pool)
	orders := postgres.NewServiceOrderRepository(pool)
	workflows := postgres.NewWorkflowRepository(pool)
	publisher := newPublisher()
	grpcServer := grpc.NewServer()
	grpcserver.New(application.NewServiceDefinitionService(services), application.NewServiceOrderService(orders, services, publisher), application.NewWorkflowService(workflows)).Register(grpcServer)

	log.Printf("operations gRPC listening on %s", addr)
	log.Fatal(grpcServer.Serve(listener))
}

func newPublisher() ports.EventPublisher {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		return operationsevents.NewNoopPublisher()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	publisher, err := operationsnats.NewPublisher(ctx, natsURL)
	if err != nil {
		log.Printf("NATS publisher disabled: %v", err)
		return operationsevents.NewNoopPublisher()
	}
	return publisher
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
