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
	var operationsClient *grpcclient.OperationsClient
	var resourceClient *grpcclient.ResourceClient
	var reportingClient *grpcclient.ReportingClient

	identityAddr := os.Getenv("IDENTITY_GRPC_ADDR")
	organizationAddr := os.Getenv("ORGANIZATION_GRPC_ADDR")
	operationsAddr := os.Getenv("OPERATIONS_GRPC_ADDR")
	resourceAddr := os.Getenv("RESOURCE_GRPC_ADDR")
	reportingAddr := os.Getenv("REPORTING_GRPC_ADDR")
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
	if operationsAddr != "" {
		client, err := grpcclient.NewOperationsClient(context.Background(), operationsAddr)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		operationsClient = client
	}
	if resourceAddr != "" {
		client, err := grpcclient.NewResourceClient(context.Background(), resourceAddr)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		resourceClient = client
	}
	if reportingAddr != "" {
		client, err := grpcclient.NewReportingClient(context.Background(), reportingAddr)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		reportingClient = client
	}

	server := app.NewServer(httpfiber.NewRouter(authService, branchClient, identityClientForRoutes, operationsClient, resourceClient, reportingClient))

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(server.Listen(":" + port))
}
