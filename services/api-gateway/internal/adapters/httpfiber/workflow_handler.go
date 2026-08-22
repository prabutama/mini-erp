package httpfiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/isapr/mini-erp/services/api-gateway/internal/domain"
	"github.com/isapr/mini-erp/services/api-gateway/internal/ports"
)

func (h *OperationsHandler) CreateWorkflow(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	var req domain.CreateWorkflowRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "Invalid JSON body")
	}
	workflow, err := h.operations.CreateWorkflow(c.UserContext(), ports.CreateWorkflowRequest{BusinessID: authContext.BusinessID, Name: req.Name, Description: req.Description})
	if err != nil {
		return mapOperationsError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(workflow)
}
func (h *OperationsHandler) ListWorkflows(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	workflows, err := h.operations.ListWorkflows(c.UserContext(), ports.ListWorkflowsRequest{BusinessID: authContext.BusinessID})
	if err != nil {
		return mapOperationsError(c, err)
	}
	return c.JSON(workflows)
}
func (h *OperationsHandler) GetWorkflow(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	workflow, err := h.operations.GetWorkflow(c.UserContext(), ports.GetWorkflowRequest{WorkflowID: c.Params("workflow_id")})
	if err != nil {
		return mapOperationsError(c, err)
	}
	if workflow.BusinessID != authContext.BusinessID {
		return writeError(c, fiber.StatusForbidden, "FORBIDDEN", "Forbidden")
	}
	return c.JSON(workflow)
}
func (h *OperationsHandler) UpdateWorkflow(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	var req domain.UpdateWorkflowRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "Invalid JSON body")
	}
	workflow, err := h.operations.UpdateWorkflow(c.UserContext(), ports.UpdateWorkflowRequest{WorkflowID: c.Params("workflow_id"), BusinessID: authContext.BusinessID, Name: req.Name, Description: req.Description, Status: req.Status})
	if err != nil {
		return mapOperationsError(c, err)
	}
	return c.JSON(workflow)
}
func (h *OperationsHandler) CreateWorkflowStatus(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	var req domain.CreateWorkflowStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "Invalid JSON body")
	}
	item, err := h.operations.CreateWorkflowStatus(c.UserContext(), ports.CreateWorkflowStatusRequest{WorkflowID: c.Params("workflow_id"), BusinessID: authContext.BusinessID, Code: req.Code, Name: req.Name, Category: req.Category, SortOrder: req.SortOrder, IsInitial: req.IsInitial, IsTerminal: req.IsTerminal})
	if err != nil {
		return mapOperationsError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(item)
}
func (h *OperationsHandler) CreateWorkflowTransition(c *fiber.Ctx) error {
	authContext := c.Locals("auth_context").(domain.AuthContext)
	var req domain.CreateWorkflowTransitionRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "Invalid JSON body")
	}
	item, err := h.operations.CreateWorkflowTransition(c.UserContext(), ports.CreateWorkflowTransitionRequest{WorkflowID: c.Params("workflow_id"), BusinessID: authContext.BusinessID, FromStatusCode: req.FromStatusCode, ToStatusCode: req.ToStatusCode})
	if err != nil {
		return mapOperationsError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(item)
}
