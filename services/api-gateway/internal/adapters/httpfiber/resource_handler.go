package httpfiber

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/isapr/mini-erp/services/api-gateway/internal/domain"
	"github.com/isapr/mini-erp/services/api-gateway/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ResourceHandler struct {
	resources ports.ResourceClient
	branches  ports.BranchClient
	orders    ports.OperationsClient
}

func NewResourceHandler(resources ports.ResourceClient, branches ports.BranchClient) *ResourceHandler {
	return &ResourceHandler{resources: resources, branches: branches}
}

func NewResourceHandlerWithOrders(resources ports.ResourceClient, branches ports.BranchClient, orders ports.OperationsClient) *ResourceHandler {
	return &ResourceHandler{resources: resources, branches: branches, orders: orders}
}

func (h *ResourceHandler) Create(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	var req domain.CreateResourceRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "Invalid JSON body")
	}
	if err := h.validateBranch(c, authContext, req.BranchID); err != nil {
		if errors.Is(err, errResponseWritten) {
			return nil
		}
		return err
	}
	resource, err := h.resources.CreateResource(c.UserContext(), ports.CreateResourceRequest{BusinessID: authContext.BusinessID, BranchID: req.BranchID, Name: req.Name, Unit: req.Unit, Type: req.Type})
	if err != nil {
		return mapResourceError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(resource)
}

func (h *ResourceHandler) List(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	branchID := c.Query("branch_id")
	if authContext.Role == domain.RoleManager || authContext.Role == domain.RoleStaff {
		if branchID == "" && len(authContext.AssignedBranchIDs) == 1 {
			branchID = authContext.AssignedBranchIDs[0]
		}
		if branchID == "" {
			return writeError(c, fiber.StatusBadRequest, "BRANCH_REQUIRED", "branch_id required")
		}
	}
	if branchID != "" {
		if err := h.validateBranch(c, authContext, branchID); err != nil {
			if errors.Is(err, errResponseWritten) {
				return nil
			}
			return err
		}
	}
	resources, err := h.resources.ListResources(c.UserContext(), ports.ListResourcesRequest{BusinessID: authContext.BusinessID, BranchID: branchID})
	if err != nil {
		return mapResourceError(c, err)
	}
	return c.JSON(resources)
}

func (h *ResourceHandler) RecordStockMovement(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	resource, err := h.resources.GetResource(c.UserContext(), ports.GetResourceRequest{ResourceID: c.Params("resource_id")})
	if err != nil {
		return mapResourceError(c, err)
	}
	if resource.BusinessID != authContext.BusinessID || !canAccessBranch(authContext, resource.BranchID) {
		return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Forbidden")
	}
	var req domain.RecordStockMovementRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "Invalid JSON body")
	}
	movement, err := h.resources.RecordStockMovement(c.UserContext(), ports.RecordStockMovementRequest{BusinessID: authContext.BusinessID, BranchID: resource.BranchID, ResourceID: resource.ResourceID, MovementType: req.MovementType, Quantity: req.Quantity, Reason: req.Reason, ServiceOrderID: req.ServiceOrderID, ActorUserID: authContext.UserID, RequestID: authContext.RequestID})
	if err != nil {
		return mapResourceError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(movement)
}

func (h *ResourceHandler) Availability(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	resource, err := h.resources.GetResource(c.UserContext(), ports.GetResourceRequest{ResourceID: c.Params("resource_id")})
	if err != nil {
		return mapResourceError(c, err)
	}
	if resource.BusinessID != authContext.BusinessID || !canAccessBranch(authContext, resource.BranchID) {
		return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Forbidden")
	}
	availability, err := h.resources.GetResourceAvailability(c.UserContext(), ports.GetResourceAvailabilityRequest{ResourceID: resource.ResourceID})
	if err != nil {
		return mapResourceError(c, err)
	}
	return c.JSON(availability)
}

func (h *ResourceHandler) RecordResourceUsage(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	if h.orders == nil {
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", "Operations client required")
	}
	order, err := h.orders.GetServiceOrder(c.UserContext(), ports.GetServiceOrderRequest{ServiceOrderID: c.Params("order_id")})
	if err != nil {
		return mapOperationsError(c, err)
	}
	if order.BusinessID != authContext.BusinessID || !canAccessBranch(authContext, order.BranchID) {
		return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Forbidden")
	}
	var req domain.RecordResourceUsageRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "Invalid JSON body")
	}
	resource, err := h.resources.GetResource(c.UserContext(), ports.GetResourceRequest{ResourceID: req.ResourceID})
	if err != nil {
		return mapResourceError(c, err)
	}
	if resource.BusinessID != authContext.BusinessID || resource.BranchID != order.BranchID {
		return writeError(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "Validation failed")
	}
	usage, err := h.resources.RecordResourceUsage(c.UserContext(), ports.RecordResourceUsageRequest{BusinessID: authContext.BusinessID, BranchID: order.BranchID, ServiceOrderID: order.ServiceOrderID, ResourceID: resource.ResourceID, Quantity: req.Quantity, Reason: req.Reason, RecordedByUserID: authContext.UserID, RequestID: authContext.RequestID})
	if err != nil {
		return mapResourceError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(usage)
}

func (h *ResourceHandler) ListResourceUsage(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	if h.orders == nil {
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", "Operations client required")
	}
	order, err := h.orders.GetServiceOrder(c.UserContext(), ports.GetServiceOrderRequest{ServiceOrderID: c.Params("order_id")})
	if err != nil {
		return mapOperationsError(c, err)
	}
	if order.BusinessID != authContext.BusinessID || !canAccessBranch(authContext, order.BranchID) {
		return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Forbidden")
	}
	usages, err := h.resources.ListResourceUsage(c.UserContext(), ports.ListResourceUsageRequest{ServiceOrderID: order.ServiceOrderID})
	if err != nil {
		return mapResourceError(c, err)
	}
	return c.JSON(usages)
}

func (h *ResourceHandler) validateBranch(c *fiber.Ctx, authContext domain.AuthContext, branchID string) error {
	if !canAccessBranch(authContext, branchID) {
		if err := writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Forbidden"); err != nil {
			return err
		}
		return errResponseWritten
	}
	branch, err := h.branches.GetBranch(c.UserContext(), ports.GetBranchRequest{BranchID: branchID})
	if err != nil {
		return mapResourceError(c, err)
	}
	if branch.BusinessID != authContext.BusinessID {
		if err := writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Forbidden"); err != nil {
			return err
		}
		return errResponseWritten
	}
	return nil
}

func mapResourceError(c *fiber.Ctx, err error) error {
	switch status.Code(err) {
	case codes.InvalidArgument:
		return writeError(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "Validation failed")
	case codes.NotFound:
		return writeError(c, fiber.StatusNotFound, "NOT_FOUND", "Not found")
	case codes.AlreadyExists:
		return writeError(c, fiber.StatusConflict, "ALREADY_EXISTS", "Already exists")
	default:
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
}
