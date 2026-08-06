package httpfiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/isapr/mini-erp/services/api-gateway/internal/domain"
)

func requireBusinessContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authContext, ok := c.Locals("auth_context").(domain.AuthContext)
		if !ok || authContext.BusinessID == "" {
			return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Forbidden")
		}
		return c.Next()
	}
}

func requireRole(roles ...domain.Role) fiber.Handler {
	allowed := map[domain.Role]struct{}{}
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(c *fiber.Ctx) error {
		authContext, ok := c.Locals("auth_context").(domain.AuthContext)
		if !ok {
			return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Forbidden")
		}
		if _, ok := allowed[authContext.Role]; !ok {
			return writeError(c, fiber.StatusForbidden, "INSUFFICIENT_PERMISSION", "Insufficient permission")
		}
		return c.Next()
	}
}

func blockPlatformAdminOnTenantRoutes() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authContext, ok := c.Locals("auth_context").(domain.AuthContext)
		if !ok {
			return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Forbidden")
		}
		if authContext.Role == domain.RolePlatformAdmin {
			return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Forbidden")
		}
		return c.Next()
	}
}

func canAccessBranch(authContext domain.AuthContext, branchID string) bool {
	if authContext.Role == domain.RoleBusinessAdmin {
		return true
	}
	if authContext.Role != domain.RoleManager && authContext.Role != domain.RoleStaff {
		return false
	}
	for _, assignedBranchID := range authContext.AssignedBranchIDs {
		if assignedBranchID == branchID {
			return true
		}
	}
	return false
}
