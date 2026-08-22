package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/isapr/mini-erp/services/resource/internal/adapters/grpcserver"
	"github.com/isapr/mini-erp/services/resource/internal/adapters/postgres"
	"github.com/isapr/mini-erp/services/resource/internal/application"
	"google.golang.org/grpc"
)

func main() {
	databaseURL := env("RESOURCE_DATABASE_URL", "postgres://mini_erp:mini_erp@localhost:5432/resource_db?sslmode=disable")
	addr := env("RESOURCE_GRPC_ADDR", ":50054")
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
	grpcserver.New(application.NewResourceService(postgres.NewResourceRepository(pool))).Register(grpcServer)
	log.Printf("resource gRPC listening on %s", addr)
	log.Fatal(grpcServer.Serve(listener))
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
