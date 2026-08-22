package httpfiber

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/isapr/mini-erp/services/api-gateway/internal/domain"
	"github.com/isapr/mini-erp/services/api-gateway/internal/ports"
)

type operationsClientStub struct{}

func (operationsClientStub) CreateServiceDefinition(_ context.Context, req ports.CreateServiceDefinitionRequest) (domain.ServiceDefinitionResponse, error) {
	return domain.ServiceDefinitionResponse{ServiceDefinitionID: "service-1", BusinessID: req.BusinessID, Name: req.Name, Code: "inspection", Status: "active"}, nil
}

func (operationsClientStub) ListServiceDefinitions(_ context.Context, req ports.ListServiceDefinitionsRequest) (domain.ListServiceDefinitionsResponse, error) {
	return domain.ListServiceDefinitionsResponse{ServiceDefinitions: []domain.ServiceDefinitionResponse{{ServiceDefinitionID: "service-1", BusinessID: req.BusinessID, Name: "Inspection", Code: "inspection", Status: "active"}}}, nil
}

func (operationsClientStub) CreateServiceOrder(_ context.Context, req ports.CreateServiceOrderRequest) (domain.ServiceOrderResponse, error) {
	return domain.ServiceOrderResponse{ServiceOrderID: "order-1", BusinessID: req.BusinessID, BranchID: req.BranchID, ServiceDefinitionID: req.ServiceDefinitionID, Title: req.Title, Status: "open", Priority: req.Priority}, nil
}

func (operationsClientStub) ListServiceOrders(_ context.Context, req ports.ListServiceOrdersRequest) (domain.ListServiceOrdersResponse, error) {
	return domain.ListServiceOrdersResponse{ServiceOrders: []domain.ServiceOrderResponse{{ServiceOrderID: "order-1", BusinessID: req.BusinessID, BranchID: req.BranchID, ServiceDefinitionID: "service-1", Title: "Fix AC", Status: "open", Priority: "normal"}}}, nil
}

func (operationsClientStub) ServiceOrderSummary(context.Context, ports.ServiceOrderSummaryRequest) (domain.ServiceOrderSummaryResponse, error) {
	return domain.ServiceOrderSummaryResponse{Total: 1, Open: 1}, nil
}

func (operationsClientStub) GetServiceOrder(context.Context, ports.GetServiceOrderRequest) (domain.ServiceOrderResponse, error) {
	return domain.ServiceOrderResponse{ServiceOrderID: "order-1", BusinessID: "business-1", BranchID: "branch-1", ServiceDefinitionID: "service-1", Title: "Fix AC", Status: "open", Priority: "normal"}, nil
}

func (operationsClientStub) TransitionServiceOrder(_ context.Context, req ports.TransitionServiceOrderRequest) (domain.ServiceOrderResponse, error) {
	return domain.ServiceOrderResponse{ServiceOrderID: req.ServiceOrderID, BusinessID: req.BusinessID, BranchID: "branch-1", ServiceDefinitionID: "service-1", Title: "Fix AC", Status: req.Status, Priority: "normal"}, nil
}

func (operationsClientStub) AssignServiceOrder(_ context.Context, req ports.AssignServiceOrderRequest) (domain.ServiceOrderAssignmentResponse, error) {
	return domain.ServiceOrderAssignmentResponse{AssignmentID: "assignment-1", ServiceOrderID: req.ServiceOrderID, BusinessID: req.BusinessID, BranchID: "branch-1", AssignedUserID: req.AssignedUserID, AssignedByUserID: req.AssignedByUserID, Status: "active"}, nil
}

func (operationsClientStub) ListAssignedServiceOrders(_ context.Context, req ports.ListAssignedServiceOrdersRequest) (domain.ListServiceOrdersResponse, error) {
	return domain.ListServiceOrdersResponse{ServiceOrders: []domain.ServiceOrderResponse{{ServiceOrderID: "order-1", BusinessID: req.BusinessID, BranchID: req.BranchID, ServiceDefinitionID: "service-1", Title: "Fix AC", Status: "open", Priority: "normal"}}}, nil
}

func (operationsClientStub) ListServiceOrderAssignments(context.Context, ports.ListServiceOrderAssignmentsRequest) (domain.ListServiceOrderAssignmentsResponse, error) {
	return domain.ListServiceOrderAssignmentsResponse{Assignments: []domain.ServiceOrderAssignmentResponse{{AssignmentID: "assignment-1", ServiceOrderID: "order-1", BusinessID: "business-1", BranchID: "branch-1", AssignedUserID: "user-2", AssignedByUserID: "user-1", Status: "active"}}}, nil
}

func (operationsClientStub) CreateWorkflow(_ context.Context, req ports.CreateWorkflowRequest) (domain.WorkflowResponse, error) {
	return domain.WorkflowResponse{WorkflowID: "workflow-1", BusinessID: req.BusinessID, Name: req.Name, Description: req.Description, Status: "active"}, nil
}

func (operationsClientStub) ListWorkflows(_ context.Context, req ports.ListWorkflowsRequest) (domain.ListWorkflowsResponse, error) {
	return domain.ListWorkflowsResponse{Workflows: []domain.WorkflowResponse{{WorkflowID: "workflow-1", BusinessID: req.BusinessID, Name: "Default", Status: "active"}}}, nil
}

func (operationsClientStub) GetWorkflow(context.Context, ports.GetWorkflowRequest) (domain.WorkflowResponse, error) {
	return domain.WorkflowResponse{WorkflowID: "workflow-1", BusinessID: "business-1", Name: "Default", Status: "active"}, nil
}

func (operationsClientStub) UpdateWorkflow(_ context.Context, req ports.UpdateWorkflowRequest) (domain.WorkflowResponse, error) {
	return domain.WorkflowResponse{WorkflowID: req.WorkflowID, BusinessID: req.BusinessID, Name: req.Name, Description: req.Description, Status: req.Status}, nil
}

func (operationsClientStub) CreateWorkflowStatus(_ context.Context, req ports.CreateWorkflowStatusRequest) (domain.WorkflowStatusResponse, error) {
	return domain.WorkflowStatusResponse{WorkflowStatusID: "status-1", WorkflowID: req.WorkflowID, BusinessID: req.BusinessID, Code: req.Code, Name: req.Name, Category: req.Category, SortOrder: req.SortOrder, IsInitial: req.IsInitial, IsTerminal: req.IsTerminal}, nil
}

func (operationsClientStub) CreateWorkflowTransition(_ context.Context, req ports.CreateWorkflowTransitionRequest) (domain.WorkflowTransitionResponse, error) {
	return domain.WorkflowTransitionResponse{WorkflowTransitionID: "transition-1", WorkflowID: req.WorkflowID, BusinessID: req.BusinessID, FromStatusCode: req.FromStatusCode, ToStatusCode: req.ToStatusCode}, nil
}

func TestManagerCanCreateServiceOrderInAssignedBranch(t *testing.T) {
	app := NewRouter(authServiceStub{role: domain.RoleManager, businessID: "business-1", assignedBranchIDs: []string{"branch-1"}}, branchClientStub{}, operationsClientStub{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/service-orders", bytes.NewReader([]byte(`{"branch_id":"branch-1","service_definition_id":"service-1","title":"Fix AC"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
}

func TestManagerCannotCreateServiceOrderInUnassignedBranch(t *testing.T) {
	app := NewRouter(authServiceStub{role: domain.RoleManager, businessID: "business-1", assignedBranchIDs: []string{"branch-1"}}, branchClientStub{}, operationsClientStub{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/service-orders", bytes.NewReader([]byte(`{"branch_id":"branch-2","service_definition_id":"service-1","title":"Fix AC"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", resp.StatusCode)
	}
}

func TestStaffCannotCreateServiceDefinition(t *testing.T) {
	app := NewRouter(authServiceStub{role: domain.RoleStaff, businessID: "business-1", assignedBranchIDs: []string{"branch-1"}}, branchClientStub{}, operationsClientStub{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/service-definitions", bytes.NewReader([]byte(`{"name":"Inspection"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", resp.StatusCode)
	}
}

func TestStaffCanTransitionAssignedServiceOrder(t *testing.T) {
	app := NewRouter(authServiceStub{role: domain.RoleStaff, businessID: "business-1", assignedBranchIDs: []string{"branch-1"}}, branchClientStub{}, operationsClientStub{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/service-orders/order-1/transition", bytes.NewReader([]byte(`{"status":"in_progress"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

type identityClientWithUserStub struct{ identityClientStubForHTTP }

func (identityClientWithUserStub) ListUsers(context.Context, ports.ListUsersRequest) (domain.ListUsersResponse, error) {
	return domain.ListUsersResponse{Users: []domain.UserResponse{{UserID: "user-2", Email: "staff@example.test", FullName: "Staff", Status: "active"}}}, nil
}

func TestBusinessAdminCanAssignServiceOrder(t *testing.T) {
	app := NewRouter(authServiceStub{role: domain.RoleBusinessAdmin, businessID: "business-1"}, branchClientStub{}, operationsClientStub{}, identityClientWithUserStub{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/service-orders/order-1/assign", bytes.NewReader([]byte(`{"assigned_user_id":"user-2"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
}

func TestStaffCannotAssignServiceOrder(t *testing.T) {
	app := NewRouter(authServiceStub{role: domain.RoleStaff, businessID: "business-1", assignedBranchIDs: []string{"branch-1"}}, branchClientStub{}, operationsClientStub{}, identityClientWithUserStub{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/service-orders/order-1/assign", bytes.NewReader([]byte(`{"assigned_user_id":"user-2"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", resp.StatusCode)
	}
}

func TestStaffCanListMyAssignedServiceOrders(t *testing.T) {
	app := NewRouter(authServiceStub{role: domain.RoleStaff, businessID: "business-1", assignedBranchIDs: []string{"branch-1"}}, branchClientStub{}, operationsClientStub{}, identityClientWithUserStub{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/service-orders/mine", nil)
	req.Header.Set("Authorization", "Bearer token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestManagerCanListServiceOrderAssignments(t *testing.T) {
	app := NewRouter(authServiceStub{role: domain.RoleManager, businessID: "business-1", assignedBranchIDs: []string{"branch-1"}}, branchClientStub{}, operationsClientStub{}, identityClientWithUserStub{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/service-orders/order-1/assignments", nil)
	req.Header.Set("Authorization", "Bearer token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestManagerCanGetServiceOrderSummaryForAssignedBranch(t *testing.T) {
	app := NewRouter(authServiceStub{role: domain.RoleManager, businessID: "business-1", assignedBranchIDs: []string{"branch-1"}}, branchClientStub{}, operationsClientStub{}, identityClientWithUserStub{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/service-orders/summary?branch_id=branch-1", nil)
	req.Header.Set("Authorization", "Bearer token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestStaffCannotFilterServiceOrdersByOtherUser(t *testing.T) {
	app := NewRouter(authServiceStub{role: domain.RoleStaff, businessID: "business-1", assignedBranchIDs: []string{"branch-1"}}, branchClientStub{}, operationsClientStub{}, identityClientWithUserStub{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/service-orders?assigned_user_id=user-2", nil)
	req.Header.Set("Authorization", "Bearer token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", resp.StatusCode)
	}
}
