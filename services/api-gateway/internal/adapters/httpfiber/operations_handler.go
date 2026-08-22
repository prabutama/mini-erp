package httpfiber

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/isapr/mini-erp/services/api-gateway/internal/domain"
	"github.com/isapr/mini-erp/services/api-gateway/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errResponseWritten = errors.New("response written")

type OperationsHandler struct {
	operations ports.OperationsClient
	branches   ports.BranchClient
	identity   ports.IdentityClient
}

func NewOperationsHandler(operations ports.OperationsClient, branches ports.BranchClient, identity ports.IdentityClient) *OperationsHandler {
	return &OperationsHandler{operations: operations, branches: branches, identity: identity}
}

func (h *OperationsHandler) CreateServiceDefinition(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	var req domain.CreateServiceDefinitionRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "Invalid JSON body")
	}
	service, err := h.operations.CreateServiceDefinition(c.UserContext(), ports.CreateServiceDefinitionRequest{BusinessID: authContext.BusinessID, Name: req.Name, Description: req.Description})
	if err != nil {
		return mapOperationsError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(service)
}

func (h *OperationsHandler) ListServiceDefinitions(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	services, err := h.operations.ListServiceDefinitions(c.UserContext(), ports.ListServiceDefinitionsRequest{BusinessID: authContext.BusinessID})
	if err != nil {
		return mapOperationsError(c, err)
	}
	return c.JSON(services)
}

func (h *OperationsHandler) CreateServiceOrder(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	var req domain.CreateServiceOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "Invalid JSON body")
	}
	if err := h.validateBranchAccess(c, authContext, req.BranchID); err != nil {
		if errors.Is(err, errResponseWritten) {
			return nil
		}
		return err
	}
	order, err := h.operations.CreateServiceOrder(c.UserContext(), ports.CreateServiceOrderRequest{BusinessID: authContext.BusinessID, BranchID: req.BranchID, ServiceDefinitionID: req.ServiceDefinitionID, Title: req.Title, Description: req.Description, Priority: req.Priority})
	if err != nil {
		return mapOperationsError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(order)
}

func (h *OperationsHandler) ListServiceOrders(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	branchID := c.Query("branch_id")
	assignedUserID := c.Query("assigned_user_id")
	if authContext.Role == domain.RoleManager || authContext.Role == domain.RoleStaff {
		if authContext.Role == domain.RoleStaff && assignedUserID != "" && assignedUserID != authContext.UserID {
			return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Forbidden")
		}
		if branchID != "" {
			if err := h.validateBranchAccess(c, authContext, branchID); err != nil {
				if errors.Is(err, errResponseWritten) {
					return nil
				}
				return err
			}
		} else if len(authContext.AssignedBranchIDs) == 1 {
			branchID = authContext.AssignedBranchIDs[0]
		} else if len(authContext.AssignedBranchIDs) == 0 {
			return c.JSON(domain.ListServiceOrdersResponse{ServiceOrders: []domain.ServiceOrderResponse{}})
		} else {
			return writeError(c, fiber.StatusBadRequest, "BRANCH_REQUIRED", "branch_id required")
		}
	}
	orders, err := h.operations.ListServiceOrders(c.UserContext(), ports.ListServiceOrdersRequest{BusinessID: authContext.BusinessID, BranchID: branchID, Status: c.Query("status"), AssignedUserID: assignedUserID})
	if err != nil {
		return mapOperationsError(c, err)
	}
	return c.JSON(orders)
}

func (h *OperationsHandler) ServiceOrderSummary(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	branchID := c.Query("branch_id")
	assignedUserID := c.Query("assigned_user_id")
	if authContext.Role == domain.RoleManager || authContext.Role == domain.RoleStaff {
		if authContext.Role == domain.RoleStaff && assignedUserID != "" && assignedUserID != authContext.UserID {
			return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Forbidden")
		}
		if branchID != "" {
			if err := h.validateBranchAccess(c, authContext, branchID); err != nil {
				if errors.Is(err, errResponseWritten) {
					return nil
				}
				return err
			}
		} else if len(authContext.AssignedBranchIDs) == 1 {
			branchID = authContext.AssignedBranchIDs[0]
		} else if len(authContext.AssignedBranchIDs) == 0 {
			return c.JSON(domain.ServiceOrderSummaryResponse{})
		} else {
			return writeError(c, fiber.StatusBadRequest, "BRANCH_REQUIRED", "branch_id required")
		}
	}
	summary, err := h.operations.ServiceOrderSummary(c.UserContext(), ports.ServiceOrderSummaryRequest{BusinessID: authContext.BusinessID, BranchID: branchID, AssignedUserID: assignedUserID})
	if err != nil {
		return mapOperationsError(c, err)
	}
	return c.JSON(summary)
}

func (h *OperationsHandler) ListMyServiceOrders(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	branchID := c.Query("branch_id")
	if branchID != "" {
		if err := h.validateBranchAccess(c, authContext, branchID); err != nil {
			if errors.Is(err, errResponseWritten) {
				return nil
			}
			return err
		}
	} else if authContext.Role == domain.RoleManager || authContext.Role == domain.RoleStaff {
		if len(authContext.AssignedBranchIDs) == 1 {
			branchID = authContext.AssignedBranchIDs[0]
		} else if len(authContext.AssignedBranchIDs) == 0 {
			return c.JSON(domain.ListServiceOrdersResponse{ServiceOrders: []domain.ServiceOrderResponse{}})
		} else {
			return writeError(c, fiber.StatusBadRequest, "BRANCH_REQUIRED", "branch_id required")
		}
	}
	orders, err := h.operations.ListAssignedServiceOrders(c.UserContext(), ports.ListAssignedServiceOrdersRequest{BusinessID: authContext.BusinessID, AssignedUserID: authContext.UserID, BranchID: branchID})
	if err != nil {
		return mapOperationsError(c, err)
	}
	return c.JSON(orders)
}

func (h *OperationsHandler) GetServiceOrder(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	order, err := h.operations.GetServiceOrder(c.UserContext(), ports.GetServiceOrderRequest{ServiceOrderID: c.Params("service_order_id")})
	if err != nil {
		return mapOperationsError(c, err)
	}
	if order.BusinessID != authContext.BusinessID || !canAccessBranch(authContext, order.BranchID) {
		return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Forbidden")
	}
	return c.JSON(order)
}

func (h *OperationsHandler) TransitionServiceOrder(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	order, err := h.operations.GetServiceOrder(c.UserContext(), ports.GetServiceOrderRequest{ServiceOrderID: c.Params("service_order_id")})
	if err != nil {
		return mapOperationsError(c, err)
	}
	if order.BusinessID != authContext.BusinessID || !canAccessBranch(authContext, order.BranchID) {
		return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Forbidden")
	}

	var req domain.TransitionServiceOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "Invalid JSON body")
	}
	updatedOrder, err := h.operations.TransitionServiceOrder(c.UserContext(), ports.TransitionServiceOrderRequest{ServiceOrderID: order.ServiceOrderID, BusinessID: authContext.BusinessID, Status: req.Status, ChangedByUserID: authContext.UserID, RequestID: authContext.RequestID})
	if err != nil {
		return mapOperationsError(c, err)
	}
	return c.JSON(updatedOrder)
}

func (h *OperationsHandler) AssignServiceOrder(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	order, err := h.operations.GetServiceOrder(c.UserContext(), ports.GetServiceOrderRequest{ServiceOrderID: c.Params("service_order_id")})
	if err != nil {
		return mapOperationsError(c, err)
	}
	if order.BusinessID != authContext.BusinessID || !canAccessBranch(authContext, order.BranchID) {
		return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Forbidden")
	}

	var req domain.AssignServiceOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "Invalid JSON body")
	}
	if h.identity == nil {
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", "Identity client required")
	}
	users, err := h.identity.ListUsers(c.UserContext(), ports.ListUsersRequest{BusinessID: authContext.BusinessID})
	if err != nil {
		return mapUserError(c, err)
	}
	if !userInBusiness(users.Users, req.AssignedUserID) {
		return writeError(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "Validation failed")
	}

	assignment, err := h.operations.AssignServiceOrder(c.UserContext(), ports.AssignServiceOrderRequest{ServiceOrderID: order.ServiceOrderID, BusinessID: authContext.BusinessID, AssignedUserID: req.AssignedUserID, AssignedByUserID: authContext.UserID, RequestID: authContext.RequestID})
	if err != nil {
		return mapOperationsError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(assignment)
}

func (h *OperationsHandler) ListServiceOrderAssignments(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	order, err := h.operations.GetServiceOrder(c.UserContext(), ports.GetServiceOrderRequest{ServiceOrderID: c.Params("service_order_id")})
	if err != nil {
		return mapOperationsError(c, err)
	}
	if order.BusinessID != authContext.BusinessID || !canAccessBranch(authContext, order.BranchID) {
		return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Forbidden")
	}
	assignments, err := h.operations.ListServiceOrderAssignments(c.UserContext(), ports.ListServiceOrderAssignmentsRequest{ServiceOrderID: order.ServiceOrderID})
	if err != nil {
		return mapOperationsError(c, err)
	}
	return c.JSON(assignments)
}

func userInBusiness(users []domain.UserResponse, userID string) bool {
	for _, user := range users {
		if user.UserID == userID {
			return true
		}
	}
	return false
}

func (h *OperationsHandler) validateBranchAccess(c *fiber.Ctx, authContext domain.AuthContext, branchID string) error {
	if !canAccessBranch(authContext, branchID) {
		if err := writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Forbidden"); err != nil {
			return err
		}
		return errResponseWritten
	}
	branch, err := h.branches.GetBranch(c.UserContext(), ports.GetBranchRequest{BranchID: branchID})
	if err != nil {
		return mapOperationsError(c, err)
	}
	if branch.BusinessID != authContext.BusinessID {
		if err := writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Forbidden"); err != nil {
			return err
		}
		return errResponseWritten
	}
	return nil
}

func mapOperationsError(c *fiber.Ctx, err error) error {
	switch status.Code(err) {
	case codes.InvalidArgument:
		return writeError(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "Validation failed")
	case codes.NotFound:
		return writeError(c, fiber.StatusNotFound, "NOT_FOUND", "Not found")
	case codes.AlreadyExists:
		return writeError(c, fiber.StatusConflict, "ALREADY_EXISTS", "Already exists")
	case codes.FailedPrecondition:
		if strings.Contains(err.Error(), "service order closed") {
			return writeError(c, fiber.StatusConflict, "SERVICE_ORDER_CLOSED", "Service order closed")
		}
		return writeError(c, fiber.StatusConflict, "INVALID_STATUS_TRANSITION", "Invalid status transition")
	default:
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
}
