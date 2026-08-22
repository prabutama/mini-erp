package application

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/isapr/mini-erp/services/operations/internal/domain"
	"github.com/isapr/mini-erp/services/operations/internal/ports"
)

type CreateWorkflowInput struct{ BusinessID, Name, Description string }
type UpdateWorkflowInput struct{ WorkflowID, BusinessID, Name, Description, Status string }
type CreateWorkflowStatusInput struct {
	WorkflowID, BusinessID, Code, Name, Category string
	SortOrder                                    int
	IsInitial, IsTerminal                        bool
}
type CreateWorkflowTransitionInput struct{ WorkflowID, BusinessID, FromStatusCode, ToStatusCode string }

type WorkflowService struct{ workflows ports.WorkflowRepository }

func NewWorkflowService(workflows ports.WorkflowRepository) *WorkflowService {
	return &WorkflowService{workflows: workflows}
}

func (s *WorkflowService) Create(ctx context.Context, input CreateWorkflowInput) (domain.Workflow, error) {
	businessID, err := uuid.Parse(input.BusinessID)
	if err != nil {
		return domain.Workflow{}, ErrValidation
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return domain.Workflow{}, ErrValidation
	}
	workflow := domain.Workflow{ID: uuid.New(), BusinessID: businessID, Name: name, Description: strings.TrimSpace(input.Description), Status: "active"}
	if err := s.workflows.Create(ctx, workflow); err != nil {
		return domain.Workflow{}, err
	}
	return workflow, nil
}

func (s *WorkflowService) List(ctx context.Context, businessID string) ([]domain.Workflow, error) {
	if _, err := uuid.Parse(businessID); err != nil {
		return nil, ErrValidation
	}
	return s.workflows.ListByBusiness(ctx, businessID)
}

func (s *WorkflowService) Get(ctx context.Context, workflowID string) (domain.Workflow, error) {
	if _, err := uuid.Parse(workflowID); err != nil {
		return domain.Workflow{}, ErrValidation
	}
	return s.workflows.FindByID(ctx, workflowID)
}

func (s *WorkflowService) Update(ctx context.Context, input UpdateWorkflowInput) (domain.Workflow, error) {
	businessID, err := uuid.Parse(input.BusinessID)
	if err != nil {
		return domain.Workflow{}, ErrValidation
	}
	workflow, err := s.Get(ctx, input.WorkflowID)
	if err != nil {
		return domain.Workflow{}, err
	}
	if workflow.BusinessID != businessID {
		return domain.Workflow{}, ErrValidation
	}
	if strings.TrimSpace(input.Name) != "" {
		workflow.Name = strings.TrimSpace(input.Name)
	}
	workflow.Description = strings.TrimSpace(input.Description)
	if strings.TrimSpace(input.Status) != "" {
		workflow.Status = strings.TrimSpace(input.Status)
	}
	if err := s.workflows.Update(ctx, workflow); err != nil {
		return domain.Workflow{}, err
	}
	return s.workflows.FindByID(ctx, input.WorkflowID)
}

func (s *WorkflowService) AddStatus(ctx context.Context, input CreateWorkflowStatusInput) (domain.WorkflowStatus, error) {
	businessID, err := uuid.Parse(input.BusinessID)
	if err != nil {
		return domain.WorkflowStatus{}, ErrValidation
	}
	workflow, err := s.Get(ctx, input.WorkflowID)
	if err != nil {
		return domain.WorkflowStatus{}, err
	}
	if workflow.BusinessID != businessID || strings.TrimSpace(input.Code) == "" || strings.TrimSpace(input.Name) == "" {
		return domain.WorkflowStatus{}, ErrValidation
	}
	category := strings.TrimSpace(input.Category)
	if category == "" {
		category = "normal"
	}
	status := domain.WorkflowStatus{ID: uuid.New(), WorkflowID: workflow.ID, BusinessID: businessID, Code: strings.TrimSpace(input.Code), Name: strings.TrimSpace(input.Name), Category: category, SortOrder: input.SortOrder, IsInitial: input.IsInitial, IsTerminal: input.IsTerminal}
	if err := s.workflows.CreateStatus(ctx, status); err != nil {
		return domain.WorkflowStatus{}, err
	}
	return status, nil
}

func (s *WorkflowService) AddTransition(ctx context.Context, input CreateWorkflowTransitionInput) (domain.WorkflowTransition, error) {
	businessID, err := uuid.Parse(input.BusinessID)
	if err != nil {
		return domain.WorkflowTransition{}, ErrValidation
	}
	workflow, err := s.Get(ctx, input.WorkflowID)
	if err != nil {
		return domain.WorkflowTransition{}, err
	}
	from, to := strings.TrimSpace(input.FromStatusCode), strings.TrimSpace(input.ToStatusCode)
	if workflow.BusinessID != businessID || from == "" || to == "" || from == to {
		return domain.WorkflowTransition{}, ErrValidation
	}
	transition := domain.WorkflowTransition{ID: uuid.New(), WorkflowID: workflow.ID, BusinessID: businessID, FromStatusCode: from, ToStatusCode: to}
	if err := s.workflows.CreateTransition(ctx, transition); err != nil {
		return domain.WorkflowTransition{}, err
	}
	return transition, nil
}

func (s *WorkflowService) ListStatuses(ctx context.Context, workflowID string) ([]domain.WorkflowStatus, error) {
	return s.workflows.ListStatuses(ctx, workflowID)
}
func (s *WorkflowService) ListTransitions(ctx context.Context, workflowID string) ([]domain.WorkflowTransition, error) {
	return s.workflows.ListTransitions(ctx, workflowID)
}
