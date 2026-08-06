package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/isapr/mini-erp/services/organization/internal/adapters/grpcserver"
	"github.com/isapr/mini-erp/services/organization/internal/adapters/postgres"
	"github.com/isapr/mini-erp/services/organization/internal/application"
	"google.golang.org/grpc"
)

func main() {
	databaseURL := env("ORGANIZATION_DATABASE_URL", "postgres://mini_erp:mini_erp@localhost:5432/organization_db?sslmode=disable")
	addr := env("ORGANIZATION_GRPC_ADDR", ":50052")

	pool, err := postgres.OpenPool(context.Background(), databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	grpcserver.New(
		application.NewBusinessService(postgres.NewBusinessRepository(pool)),
		application.NewBranchService(postgres.NewBranchRepository(pool)),
		application.NewPlacementService(postgres.NewPlacementRepository(pool)),
	).Register(grpcServer)

	log.Printf("organization gRPC listening on %s", addr)
	log.Fatal(grpcServer.Serve(listener))
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
