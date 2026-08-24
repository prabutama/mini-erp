package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/isapr/mini-erp/services/reporting/internal/adapters/grpcserver"
	reportingnats "github.com/isapr/mini-erp/services/reporting/internal/adapters/nats"
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
	reportingService := application.NewReportingService(postgres.NewReportingRepository(pool))
	if consumer := newConsumer(reportingService); consumer != nil {
		defer consumer.Close()
	}
	grpcServer := grpc.NewServer()
	grpcserver.New(reportingService).Register(grpcServer)
	log.Printf("reporting gRPC listening on %s", addr)
	log.Fatal(grpcServer.Serve(listener))
}

func newConsumer(service *application.ReportingService) *reportingnats.Consumer {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	consumer, err := reportingnats.NewConsumer(ctx, natsURL, service)
	if err != nil {
		log.Printf("NATS consumer disabled: %v", err)
		return nil
	}
	if err := consumer.Start(); err != nil {
		log.Printf("NATS consumer disabled: %v", err)
		consumer.Close()
		return nil
	}
	log.Printf("reporting NATS consumer started")
	return consumer
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
