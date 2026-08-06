package main

import (
	"context"
	"log"
	"os"

	"github.com/isapr/mini-erp/services/api-gateway/internal/adapters/grpcclient"
	"github.com/isapr/mini-erp/services/api-gateway/internal/adapters/httpfiber"
	"github.com/isapr/mini-erp/services/api-gateway/internal/app"
	"github.com/isapr/mini-erp/services/api-gateway/internal/service"
)

func main() {
	authService := service.NewAuthService()
	var branchClient *grpcclient.OrganizationClient
	var identityClientForRoutes *grpcclient.IdentityClient

	identityAddr := os.Getenv("IDENTITY_GRPC_ADDR")
	organizationAddr := os.Getenv("ORGANIZATION_GRPC_ADDR")
	if identityAddr != "" && organizationAddr != "" {
		identityClient, err := grpcclient.NewIdentityClient(context.Background(), identityAddr)
		if err != nil {
			log.Fatal(err)
		}
		defer identityClient.Close()
		identityClientForRoutes = identityClient

		organizationClient, err := grpcclient.NewOrganizationClient(context.Background(), organizationAddr)
		if err != nil {
			log.Fatal(err)
		}
		defer organizationClient.Close()
		branchClient = organizationClient

		authService = service.NewAuthServiceWithClients(identityClient, organizationClient)
	}

	server := app.NewServer(httpfiber.NewRouter(authService, branchClient, identityClientForRoutes))

	log.Fatal(server.Listen(":8080"))
}
