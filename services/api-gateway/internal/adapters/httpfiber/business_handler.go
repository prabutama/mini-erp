package httpfiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/isapr/mini-erp/services/api-gateway/internal/domain"
	"github.com/isapr/mini-erp/services/api-gateway/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BusinessHandler struct {
	organization ports.OrganizationClient
}

func NewBusinessHandler(organization ports.OrganizationClient) *BusinessHandler {
	return &BusinessHandler{organization: organization}
}

func (h *BusinessHandler) Current(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	business, err := h.organization.GetBusiness(c.UserContext(), ports.GetBusinessRequest{BusinessID: authContext.BusinessID})
	if err != nil {
		return mapBusinessError(c, err)
	}
	return c.JSON(business)
}

func (h *BusinessHandler) UpdateCurrent(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	var req domain.UpdateBusinessRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "Invalid JSON body")
	}
	business, err := h.organization.UpdateBusiness(c.UserContext(), ports.UpdateBusinessRequest{BusinessID: authContext.BusinessID, Name: req.Name, Timezone: req.Timezone})
	if err != nil {
		return mapBusinessError(c, err)
	}
	return c.JSON(business)
}

func (h *BusinessHandler) ListPlatform(c *fiber.Ctx) error {
	businesses, err := h.organization.ListPlatformBusinesses(c.UserContext(), ports.ListPlatformBusinessesRequest{})
	if err != nil {
		return mapBusinessError(c, err)
	}
	return c.JSON(businesses)
}

func (h *BusinessHandler) GetPlatform(c *fiber.Ctx) error {
	business, err := h.organization.GetBusiness(c.UserContext(), ports.GetBusinessRequest{BusinessID: c.Params("business_id")})
	if err != nil {
		return mapBusinessError(c, err)
	}
	return c.JSON(business)
}

func (h *BusinessHandler) UpdatePlatform(c *fiber.Ctx) error {
	var req domain.UpdatePlatformBusinessRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "Invalid JSON body")
	}
	business, err := h.organization.UpdatePlatformBusiness(c.UserContext(), ports.UpdatePlatformBusinessRequest{BusinessID: c.Params("business_id"), Status: req.Status, Plan: req.Plan, PlatformNotes: req.PlatformNotes, SuspendedAt: req.SuspendedAt})
	if err != nil {
		return mapBusinessError(c, err)
	}
	return c.JSON(business)
}

func mapBusinessError(c *fiber.Ctx, err error) error {
	switch status.Code(err) {
	case codes.InvalidArgument:
		return writeError(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "Validation failed")
	case codes.NotFound:
		return writeError(c, fiber.StatusNotFound, "BUSINESS_NOT_FOUND", "Business not found")
	default:
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
}
