package ports

import (
	"context"

	"github.com/isapr/mini-erp/services/api-gateway/internal/domain"
)

type OperationsClient interface {
	CreateServiceDefinition(ctx context.Context, req CreateServiceDefinitionRequest) (domain.ServiceDefinitionResponse, error)
	ListServiceDefinitions(ctx context.Context, req ListServiceDefinitionsRequest) (domain.ListServiceDefinitionsResponse, error)
	CreateServiceOrder(ctx context.Context, req CreateServiceOrderRequest) (domain.ServiceOrderResponse, error)
	ListServiceOrders(ctx context.Context, req ListServiceOrdersRequest) (domain.ListServiceOrdersResponse, error)
	ServiceOrderSummary(ctx context.Context, req ServiceOrderSummaryRequest) (domain.ServiceOrderSummaryResponse, error)
	GetServiceOrder(ctx context.Context, req GetServiceOrderRequest) (domain.ServiceOrderResponse, error)
	TransitionServiceOrder(ctx context.Context, req TransitionServiceOrderRequest) (domain.ServiceOrderResponse, error)
	AssignServiceOrder(ctx context.Context, req AssignServiceOrderRequest) (domain.ServiceOrderAssignmentResponse, error)
	ListAssignedServiceOrders(ctx context.Context, req ListAssignedServiceOrdersRequest) (domain.ListServiceOrdersResponse, error)
	ListServiceOrderAssignments(ctx context.Context, req ListServiceOrderAssignmentsRequest) (domain.ListServiceOrderAssignmentsResponse, error)
	CreateWorkflow(ctx context.Context, req CreateWorkflowRequest) (domain.WorkflowResponse, error)
	ListWorkflows(ctx context.Context, req ListWorkflowsRequest) (domain.ListWorkflowsResponse, error)
	GetWorkflow(ctx context.Context, req GetWorkflowRequest) (domain.WorkflowResponse, error)
	UpdateWorkflow(ctx context.Context, req UpdateWorkflowRequest) (domain.WorkflowResponse, error)
	CreateWorkflowStatus(ctx context.Context, req CreateWorkflowStatusRequest) (domain.WorkflowStatusResponse, error)
	CreateWorkflowTransition(ctx context.Context, req CreateWorkflowTransitionRequest) (domain.WorkflowTransitionResponse, error)
}

type CreateServiceDefinitionRequest struct {
	BusinessID  string `json:"business_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ListServiceDefinitionsRequest struct {
	BusinessID string `json:"business_id"`
}

type CreateServiceOrderRequest struct {
	BusinessID          string `json:"business_id"`
	BranchID            string `json:"branch_id"`
	ServiceDefinitionID string `json:"service_definition_id"`
	Title               string `json:"title"`
	Description         string `json:"description"`
	Priority            string `json:"priority"`
}

type ListServiceOrdersRequest struct {
	BusinessID     string `json:"business_id"`
	BranchID       string `json:"branch_id"`
	Status         string `json:"status"`
	AssignedUserID string `json:"assigned_user_id"`
}

type ServiceOrderSummaryRequest struct {
	BusinessID     string `json:"business_id"`
	BranchID       string `json:"branch_id"`
	AssignedUserID string `json:"assigned_user_id"`
}

type GetServiceOrderRequest struct {
	ServiceOrderID string `json:"service_order_id"`
}

type TransitionServiceOrderRequest struct {
	ServiceOrderID  string `json:"service_order_id"`
	BusinessID      string `json:"business_id"`
	Status          string `json:"status"`
	ChangedByUserID string `json:"changed_by_user_id"`
	RequestID       string `json:"request_id"`
}

type AssignServiceOrderRequest struct {
	ServiceOrderID   string `json:"service_order_id"`
	BusinessID       string `json:"business_id"`
	AssignedUserID   string `json:"assigned_user_id"`
	AssignedByUserID string `json:"assigned_by_user_id"`
	RequestID        string `json:"request_id"`
}

type ListAssignedServiceOrdersRequest struct {
	BusinessID     string `json:"business_id"`
	AssignedUserID string `json:"assigned_user_id"`
	BranchID       string `json:"branch_id"`
}

type ListServiceOrderAssignmentsRequest struct {
	ServiceOrderID string `json:"service_order_id"`
}

type CreateWorkflowRequest struct {
	BusinessID  string `json:"business_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
type ListWorkflowsRequest struct {
	BusinessID string `json:"business_id"`
}
type GetWorkflowRequest struct {
	WorkflowID string `json:"workflow_id"`
}
type UpdateWorkflowRequest struct {
	WorkflowID  string `json:"workflow_id"`
	BusinessID  string `json:"business_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
type CreateWorkflowStatusRequest struct {
	WorkflowID string `json:"workflow_id"`
	BusinessID string `json:"business_id"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	Category   string `json:"category"`
	SortOrder  int    `json:"sort_order"`
	IsInitial  bool   `json:"is_initial"`
	IsTerminal bool   `json:"is_terminal"`
}
type CreateWorkflowTransitionRequest struct {
	WorkflowID     string `json:"workflow_id"`
	BusinessID     string `json:"business_id"`
	FromStatusCode string `json:"from_status_code"`
	ToStatusCode   string `json:"to_status_code"`
}
