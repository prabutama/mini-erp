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

	v1 := app.Group("/api/v1")
	authGroup := v1.Group("/auth")
	authGroup.Post("/signup", handler.Signup)
	authGroup.Post("/login", handler.Login)
	authGroup.Post("/refresh", handler.Refresh)
	authGroup.Post("/logout", handler.Logout)
	v1.Get("/me", authMiddleware(auth), handler.Me)

	var branchClient ports.BranchClient
	var identityClient ports.IdentityClient
	for _, client := range clients {
		if c, ok := client.(ports.BranchClient); ok {
			branchClient = c
		}
		if c, ok := client.(ports.IdentityClient); ok {
			identityClient = c
		}
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
	}

	return app
}
