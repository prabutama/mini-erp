package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/isapr/mini-erp/services/reporting/internal/adapters/grpcserver"
	"github.com/isapr/mini-erp/services/reporting/internal/adapters/postgres"
	"github.com/isapr/mini-erp/services/reporting/internal/application"
	"google.golang.org/grpc"
)

func main() {
	databaseURL := env("REPORTING_DATABASE_URL", "postgres://mini_erp:mini_erp@localhost:5432/reporting_db?sslmode=disable")
	addr := env("REPORTING_GRPC_ADDR", ":50055")
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
	grpcserver.New(application.NewReportingService(postgres.NewReportingRepository(pool))).Register(grpcServer)
	log.Printf("reporting gRPC listening on %s", addr)
	log.Fatal(grpcServer.Serve(listener))
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
