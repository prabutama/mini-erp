package httpfiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/isapr/mini-erp/services/api-gateway/internal/domain"
	"github.com/isapr/mini-erp/services/api-gateway/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	identity     ports.IdentityClient
	organization ports.BranchClient
}

func NewUserHandler(identity ports.IdentityClient, organization ports.BranchClient) *UserHandler {
	return &UserHandler{identity: identity, organization: organization}
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	users, err := h.identity.ListUsers(c.UserContext(), ports.ListUsersRequest{BusinessID: authContext.BusinessID})
	if err != nil {
		return mapUserError(c, err)
	}
	return c.JSON(users)
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	var req domain.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "Invalid JSON body")
	}
	user, err := h.identity.CreateUser(c.UserContext(), ports.CreateUserRequest{BusinessID: authContext.BusinessID, Email: req.Email, Password: req.Password, FullName: req.FullName, Role: req.Role})
	if err != nil {
		return mapUserError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) Get(c *fiber.Ctx) error {
	user, err := h.identity.GetUser(c.UserContext(), ports.GetUserRequest{UserID: c.Params("user_id")})
	if err != nil {
		return mapUserError(c, err)
	}
	return c.JSON(user)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	var req domain.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "Invalid JSON body")
	}
	user, err := h.identity.UpdateUser(c.UserContext(), ports.UpdateUserRequest{UserID: c.Params("user_id"), FullName: req.FullName, Status: req.Status})
	if err != nil {
		return mapUserError(c, err)
	}
	return c.JSON(user)
}

func (h *UserHandler) AssignRole(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	var req domain.AssignRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "Invalid JSON body")
	}
	if err := h.identity.AssignBusinessRole(c.UserContext(), ports.AssignBusinessRoleRequest{UserID: c.Params("user_id"), BusinessID: authContext.BusinessID, Role: req.Role}); err != nil {
		return mapUserError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *UserHandler) CreatePlacement(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	var req domain.CreatePlacementRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "Invalid JSON body")
	}
	branch, err := h.organization.GetBranch(c.UserContext(), ports.GetBranchRequest{BranchID: req.BranchID})
	if err != nil {
		return mapUserError(c, err)
	}
	if branch.BusinessID != authContext.BusinessID {
		return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Forbidden")
	}
	placement, err := h.organization.CreateEmployeePlacement(c.UserContext(), ports.CreateEmployeePlacementRequest{UserID: c.Params("user_id"), BusinessID: authContext.BusinessID, BranchID: req.BranchID, Position: req.Position, EmploymentType: req.EmploymentType})
	if err != nil {
		return mapUserError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(placement)
}

func mapUserError(c *fiber.Ctx, err error) error {
	switch status.Code(err) {
	case codes.InvalidArgument:
		return writeError(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "Validation failed")
	case codes.AlreadyExists:
		return writeError(c, fiber.StatusConflict, "USER_ALREADY_EXISTS", "User already exists")
	case codes.NotFound:
		return writeError(c, fiber.StatusNotFound, "USER_NOT_FOUND", "User not found")
	default:
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
}
