package httpfiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/isapr/mini-erp/services/api-gateway/internal/domain"
)

func Roles(c *fiber.Ctx) error {
	return c.JSON(domain.ListRolesResponse{Roles: []domain.RoleResponse{
		{Name: domain.RolePlatformAdmin, Scope: "platform"},
		{Name: domain.RoleBusinessAdmin, Scope: "business"},
		{Name: domain.RoleManager, Scope: "branch"},
		{Name: domain.RoleStaff, Scope: "branch"},
	}})
}
