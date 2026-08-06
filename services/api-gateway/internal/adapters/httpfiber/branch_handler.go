package httpfiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/isapr/mini-erp/services/api-gateway/internal/domain"
	"github.com/isapr/mini-erp/services/api-gateway/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BranchHandler struct {
	branches ports.BranchClient
}

func NewBranchHandler(branches ports.BranchClient) *BranchHandler {
	return &BranchHandler{branches: branches}
}

func (h *BranchHandler) Create(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	if authContext.BusinessID == "" {
		return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Business context required")
	}

	var req domain.CreateBranchRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "Invalid JSON body")
	}

	branch, err := h.branches.CreateBranch(c.UserContext(), ports.CreateBranchRequest{BusinessID: authContext.BusinessID, Name: req.Name, Address: req.Address, Phone: req.Phone})
	if err != nil {
		return mapBranchError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(branch)
}

func (h *BranchHandler) List(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	if authContext.BusinessID == "" {
		return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Business context required")
	}
	branches, err := h.branches.ListBranches(c.UserContext(), ports.ListBranchesRequest{BusinessID: authContext.BusinessID})
	if err != nil {
		return mapBranchError(c, err)
	}
	if authContext.Role == domain.RoleManager || authContext.Role == domain.RoleStaff {
		filtered := domain.ListBranchesResponse{Branches: []domain.BranchResponse{}}
		for _, branch := range branches.Branches {
			if canAccessBranch(authContext, branch.BranchID) {
				filtered.Branches = append(filtered.Branches, branch)
			}
		}
		return c.JSON(filtered)
	}
	return c.JSON(branches)
}

func (h *BranchHandler) Get(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	branchID := c.Params("branch_id")
	if !canAccessBranch(authContext, branchID) {
		return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Forbidden")
	}
	branch, err := h.branches.GetBranch(c.UserContext(), ports.GetBranchRequest{BranchID: branchID})
	if err != nil {
		return mapBranchError(c, err)
	}
	if branch.BusinessID != authContext.BusinessID {
		return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Forbidden")
	}
	return c.JSON(branch)
}

func (h *BranchHandler) Update(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	branchID := c.Params("branch_id")
	existingBranch, err := h.branches.GetBranch(c.UserContext(), ports.GetBranchRequest{BranchID: branchID})
	if err != nil {
		return mapBranchError(c, err)
	}
	if existingBranch.BusinessID != authContext.BusinessID {
		return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Forbidden")
	}

	var req domain.UpdateBranchRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "Invalid JSON body")
	}
	branch, err := h.branches.UpdateBranch(c.UserContext(), ports.UpdateBranchRequest{BranchID: branchID, Name: req.Name, Address: req.Address, Phone: req.Phone, Status: req.Status})
	if err != nil {
		return mapBranchError(c, err)
	}
	return c.JSON(branch)
}

func mapBranchError(c *fiber.Ctx, err error) error {
	switch status.Code(err) {
	case codes.InvalidArgument:
		return writeError(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "Validation failed")
	case codes.NotFound:
		return writeError(c, fiber.StatusNotFound, "BRANCH_NOT_FOUND", "Branch not found")
	default:
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
}
