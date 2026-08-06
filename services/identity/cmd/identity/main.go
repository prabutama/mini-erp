package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/isapr/mini-erp/services/identity/internal/adapters/grpcserver"
	"github.com/isapr/mini-erp/services/identity/internal/adapters/postgres"
	"github.com/isapr/mini-erp/services/identity/internal/application"
	"google.golang.org/grpc"
)

func main() {
	databaseURL := env("IDENTITY_DATABASE_URL", "postgres://mini_erp:mini_erp@localhost:5432/identity_db?sslmode=disable")
	addr := env("IDENTITY_GRPC_ADDR", ":50051")

	pool, err := postgres.OpenPool(context.Background(), databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	users := postgres.NewUserRepository(pool)
	authService := application.NewAuthService(users, env("IDENTITY_JWT_SECRET", "dev-secret"))
	signupService := application.NewSignupServiceWithAuth(users, authService)

	grpcServer := grpc.NewServer()
	grpcserver.New(signupService, authService, application.NewUserService(users)).Register(grpcServer)

	log.Printf("identity gRPC listening on %s", addr)
	log.Fatal(grpcServer.Serve(listener))
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
