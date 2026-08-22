package httpfiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/isapr/mini-erp/services/api-gateway/internal/domain"
	"github.com/isapr/mini-erp/services/api-gateway/internal/ports"
)

func NewRouter(auth ports.AuthService, clients ...any) *fiber.App {
	app := fiber.New()
	handler := NewAuthHandler(auth)

	app.Use(requestIDMiddleware())
	app.Use(requestLogMiddleware())
	app.Get("/healthz", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"status": "ok"}) })
	app.Get("/readyz", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"status": "ready"}) })

	v1 := app.Group("/api/v1")
	authGroup := v1.Group("/auth")
	authGroup.Post("/signup", handler.Signup)
	authGroup.Post("/login", handler.Login)
	authGroup.Post("/refresh", handler.Refresh)
	authGroup.Post("/logout", handler.Logout)
	v1.Get("/me", authMiddleware(auth), handler.Me)

	var branchClient ports.BranchClient
	var identityClient ports.IdentityClient
	var operationsClient ports.OperationsClient
	var resourceClient ports.ResourceClient
	var reportingClient ports.ReportingClient
	for _, client := range clients {
		if c, ok := client.(ports.BranchClient); ok {
			branchClient = c
		}
		if c, ok := client.(ports.IdentityClient); ok {
			identityClient = c
		}
		if c, ok := client.(ports.OperationsClient); ok {
			operationsClient = c
		}
		if c, ok := client.(ports.ResourceClient); ok {
			resourceClient = c
		}
		if c, ok := client.(ports.ReportingClient); ok {
			reportingClient = c
		}
	}

	if identityClient != nil || branchClient != nil || operationsClient != nil || resourceClient != nil || reportingClient != nil {
		v1.Get("/roles", authMiddleware(auth), Roles)
	}

	if branchClient != nil {
		branchHandler := NewBranchHandler(branchClient)
		branches := v1.Group("/branches", authMiddleware(auth), blockPlatformAdminOnTenantRoutes(), requireBusinessContext())
		branches.Get("", branchHandler.List)
		branches.Post("", requireRole(domain.RoleBusinessAdmin), branchHandler.Create)
		branches.Get("/:branch_id", branchHandler.Get)
		branches.Patch("/:branch_id", requireRole(domain.RoleBusinessAdmin), branchHandler.Update)

		if identityClient != nil {
			userHandler := NewUserHandler(identityClient, branchClient)
			users := v1.Group("/users", authMiddleware(auth), blockPlatformAdminOnTenantRoutes(), requireBusinessContext(), requireRole(domain.RoleBusinessAdmin))
			users.Get("", userHandler.List)
			users.Post("", userHandler.Create)
			users.Get("/:user_id", userHandler.Get)
			users.Patch("/:user_id", userHandler.Update)
			users.Post("/:user_id/roles", userHandler.AssignRole)
			users.Post("/:user_id/placements", userHandler.CreatePlacement)
		}

		if organizationClient, ok := branchClient.(ports.OrganizationClient); ok {
			businessHandler := NewBusinessHandler(organizationClient)
			currentBusiness := v1.Group("/businesses/current", authMiddleware(auth), blockPlatformAdminOnTenantRoutes(), requireBusinessContext())
			currentBusiness.Get("", businessHandler.Current)
			currentBusiness.Patch("", requireRole(domain.RoleBusinessAdmin), businessHandler.UpdateCurrent)

			platform := v1.Group("/platform", authMiddleware(auth), requireRole(domain.RolePlatformAdmin))
			platform.Get("/businesses", businessHandler.ListPlatform)
			platform.Get("/businesses/:business_id", businessHandler.GetPlatform)
			platform.Patch("/businesses/:business_id", businessHandler.UpdatePlatform)
		}

		if operationsClient != nil {
			operationsHandler := NewOperationsHandler(operationsClient, branchClient, identityClient)
			workflows := v1.Group("/workflows", authMiddleware(auth), blockPlatformAdminOnTenantRoutes(), requireBusinessContext())
			workflows.Get("", operationsHandler.ListWorkflows)
			workflows.Post("", requireRole(domain.RoleBusinessAdmin, domain.RoleManager), operationsHandler.CreateWorkflow)
			workflows.Get("/:workflow_id", operationsHandler.GetWorkflow)
			workflows.Patch("/:workflow_id", requireRole(domain.RoleBusinessAdmin, domain.RoleManager), operationsHandler.UpdateWorkflow)
			workflows.Post("/:workflow_id/statuses", requireRole(domain.RoleBusinessAdmin, domain.RoleManager), operationsHandler.CreateWorkflowStatus)
			workflows.Post("/:workflow_id/transitions", requireRole(domain.RoleBusinessAdmin, domain.RoleManager), operationsHandler.CreateWorkflowTransition)

			serviceDefinitions := v1.Group("/service-definitions", authMiddleware(auth), blockPlatformAdminOnTenantRoutes(), requireBusinessContext())
			serviceDefinitions.Get("", operationsHandler.ListServiceDefinitions)
			serviceDefinitions.Post("", requireRole(domain.RoleBusinessAdmin, domain.RoleManager), operationsHandler.CreateServiceDefinition)

			serviceOrders := v1.Group("/service-orders", authMiddleware(auth), blockPlatformAdminOnTenantRoutes(), requireBusinessContext())
			serviceOrders.Get("", operationsHandler.ListServiceOrders)
			serviceOrders.Get("/summary", operationsHandler.ServiceOrderSummary)
			serviceOrders.Get("/mine", operationsHandler.ListMyServiceOrders)
			serviceOrders.Post("", requireRole(domain.RoleBusinessAdmin, domain.RoleManager), operationsHandler.CreateServiceOrder)
			serviceOrders.Get("/:service_order_id", operationsHandler.GetServiceOrder)
			serviceOrders.Get("/:service_order_id/assignments", operationsHandler.ListServiceOrderAssignments)
			serviceOrders.Post("/:service_order_id/assign", requireRole(domain.RoleBusinessAdmin, domain.RoleManager), operationsHandler.AssignServiceOrder)
			serviceOrders.Post("/:service_order_id/transition", operationsHandler.TransitionServiceOrder)
		}

		if resourceClient != nil {
			resourceHandler := NewResourceHandlerWithOrders(resourceClient, branchClient, operationsClient)
			resources := v1.Group("/resources", authMiddleware(auth), blockPlatformAdminOnTenantRoutes(), requireBusinessContext())
			resources.Get("", resourceHandler.List)
			resources.Post("", requireRole(domain.RoleBusinessAdmin, domain.RoleManager), resourceHandler.Create)
			resources.Post("/:resource_id/stock-movements", requireRole(domain.RoleBusinessAdmin, domain.RoleManager, domain.RoleStaff), resourceHandler.RecordStockMovement)
			resources.Get("/:resource_id/availability", resourceHandler.Availability)

			resourceServiceOrders := v1.Group("/service-orders", authMiddleware(auth), blockPlatformAdminOnTenantRoutes(), requireBusinessContext())
			resourceServiceOrders.Post("/:order_id/resource-usage", requireRole(domain.RoleBusinessAdmin, domain.RoleManager, domain.RoleStaff), resourceHandler.RecordResourceUsage)
			resourceServiceOrders.Get("/:order_id/resource-usage", resourceHandler.ListResourceUsage)
		}

		if reportingClient != nil {
			reportingHandler := NewReportingHandler(reportingClient, branchClient)
			reports := v1.Group("/reports", authMiddleware(auth), blockPlatformAdminOnTenantRoutes(), requireBusinessContext(), requireRole(domain.RoleBusinessAdmin, domain.RoleManager))
			reports.Get("/audit-events", reportingHandler.AuditEvents)
			reports.Get("/operations-summary", reportingHandler.OperationsSummary)
		}
	}

	return app
}
