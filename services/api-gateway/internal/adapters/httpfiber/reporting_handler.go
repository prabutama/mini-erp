package httpfiber

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/isapr/mini-erp/services/api-gateway/internal/domain"
	"github.com/isapr/mini-erp/services/api-gateway/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReportingHandler struct {
	reports  ports.ReportingClient
	branches ports.BranchClient
}

func NewReportingHandler(reports ports.ReportingClient, branches ports.BranchClient) *ReportingHandler {
	return &ReportingHandler{reports: reports, branches: branches}
}

func (h *ReportingHandler) AuditEvents(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	branchID, ok := h.branchFilter(c, authContext)
	if !ok {
		return nil
	}
	events, err := h.reports.GetAuditEvents(c.UserContext(), ports.GetAuditEventsRequest{BusinessID: authContext.BusinessID, BranchID: branchID})
	if err != nil {
		return mapReportingError(c, err)
	}
	return c.JSON(events)
}

func (h *ReportingHandler) OperationsSummary(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	branchID, ok := h.branchFilter(c, authContext)
	if !ok {
		return nil
	}
	summary, err := h.reports.GetOperationsSummary(c.UserContext(), ports.GetOperationsSummaryRequest{BusinessID: authContext.BusinessID, BranchID: branchID, SnapshotDate: c.Query("date")})
	if err != nil {
		return mapReportingError(c, err)
	}
	return c.JSON(summary)
}

func (h *ReportingHandler) branchFilter(c *fiber.Ctx, authContext domain.AuthContext) (string, bool) {
	if authContext.Role == domain.RoleStaff {
		_ = writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Forbidden")
		return "", false
	}
	branchID := c.Query("branch_id")
	if authContext.Role == domain.RoleManager {
		if branchID == "" && len(authContext.AssignedBranchIDs) == 1 {
			branchID = authContext.AssignedBranchIDs[0]
		}
		if branchID == "" {
			_ = writeError(c, fiber.StatusBadRequest, "BRANCH_REQUIRED", "branch_id required")
			return "", false
		}
	}
	if branchID != "" {
		if err := h.validateBranch(c, authContext, branchID); err != nil {
			return "", false
		}
	}
	return branchID, true
}

func (h *ReportingHandler) validateBranch(c *fiber.Ctx, authContext domain.AuthContext, branchID string) error {
	if !canAccessBranch(authContext, branchID) {
		return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Forbidden")
	}
	branch, err := h.branches.GetBranch(c.UserContext(), ports.GetBranchRequest{BranchID: branchID})
	if err != nil {
		return mapReportingError(c, err)
	}
	if branch.BusinessID != authContext.BusinessID {
		return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Forbidden")
	}
	return nil
}

func mapReportingError(c *fiber.Ctx, err error) error {
	switch status.Code(err) {
	case codes.InvalidArgument:
		return writeError(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "Validation failed")
	case codes.NotFound:
		return writeError(c, fiber.StatusNotFound, "NOT_FOUND", "Not found")
	default:
		if errors.Is(err, errResponseWritten) {
			return nil
		}
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
}
