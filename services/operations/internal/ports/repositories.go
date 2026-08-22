package ports

import (
	"context"

	"github.com/isapr/mini-erp/services/operations/internal/domain"
)

type ServiceDefinitionRepository interface {
	Create(ctx context.Context, service domain.ServiceDefinition) error
	ListByBusiness(ctx context.Context, businessID string) ([]domain.ServiceDefinition, error)
	FindByID(ctx context.Context, serviceDefinitionID string) (domain.ServiceDefinition, error)
}

type ServiceOrderRepository interface {
	Create(ctx context.Context, order domain.ServiceOrder) error
	ListByBusiness(ctx context.Context, businessID string, branchID string, status string, assignedUserID string) ([]domain.ServiceOrder, error)
	SummaryByBusiness(ctx context.Context, businessID string, branchID string, assignedUserID string) (domain.ServiceOrderSummary, error)
	FindByID(ctx context.Context, serviceOrderID string) (domain.ServiceOrder, error)
	UpdateStatus(ctx context.Context, order domain.ServiceOrder, previousStatus string, changedByUserID string, requestID string) error
	Assign(ctx context.Context, assignment domain.ServiceOrderAssignment) error
	ListAssignedToUser(ctx context.Context, businessID string, assignedUserID string, branchID string) ([]domain.ServiceOrder, error)
	ListAssignmentsByOrder(ctx context.Context, serviceOrderID string) ([]domain.ServiceOrderAssignment, error)
}

type WorkflowRepository interface {
	Create(ctx context.Context, workflow domain.Workflow) error
	ListByBusiness(ctx context.Context, businessID string) ([]domain.Workflow, error)
	FindByID(ctx context.Context, workflowID string) (domain.Workflow, error)
	Update(ctx context.Context, workflow domain.Workflow) error
	CreateStatus(ctx context.Context, status domain.WorkflowStatus) error
	ListStatuses(ctx context.Context, workflowID string) ([]domain.WorkflowStatus, error)
	CreateTransition(ctx context.Context, transition domain.WorkflowTransition) error
	ListTransitions(ctx context.Context, workflowID string) ([]domain.WorkflowTransition, error)
}
