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

type resourceClientStub struct{}

func (resourceClientStub) CreateResource(_ context.Context, req ports.CreateResourceRequest) (domain.ResourceResponse, error) {
	return domain.ResourceResponse{ResourceID: "resource-1", BusinessID: req.BusinessID, BranchID: req.BranchID, Name: req.Name, Code: "filter", Unit: req.Unit, Type: req.Type, Status: "active"}, nil
}
func (resourceClientStub) ListResources(_ context.Context, req ports.ListResourcesRequest) (domain.ListResourcesResponse, error) {
	return domain.ListResourcesResponse{Resources: []domain.ResourceResponse{{ResourceID: "resource-1", BusinessID: req.BusinessID, BranchID: req.BranchID, Name: "Filter", Code: "filter", Unit: "pcs", Type: "stock", Status: "active"}}}, nil
}
func (resourceClientStub) GetResource(context.Context, ports.GetResourceRequest) (domain.ResourceResponse, error) {
	return domain.ResourceResponse{ResourceID: "resource-1", BusinessID: "business-1", BranchID: "branch-1", Name: "Filter", Code: "filter", Unit: "pcs", Type: "stock", Status: "active"}, nil
}
func (resourceClientStub) RecordStockMovement(_ context.Context, req ports.RecordStockMovementRequest) (domain.StockMovementResponse, error) {
	return domain.StockMovementResponse{StockMovementID: "movement-1", BusinessID: req.BusinessID, BranchID: req.BranchID, ResourceID: req.ResourceID, MovementType: req.MovementType, Quantity: req.Quantity}, nil
}
func (resourceClientStub) GetResourceAvailability(context.Context, ports.GetResourceAvailabilityRequest) (domain.ResourceAvailabilityResponse, error) {
	return domain.ResourceAvailabilityResponse{ResourceID: "resource-1", Quantity: 5}, nil
}
func (resourceClientStub) RecordResourceUsage(_ context.Context, req ports.RecordResourceUsageRequest) (domain.ResourceUsageResponse, error) {
	return domain.ResourceUsageResponse{ResourceUsageID: "usage-1", BusinessID: req.BusinessID, BranchID: req.BranchID, ServiceOrderID: req.ServiceOrderID, ResourceID: req.ResourceID, Quantity: req.Quantity, RecordedByUserID: req.RecordedByUserID, StockMovementID: "movement-1"}, nil
}
func (resourceClientStub) ListResourceUsage(context.Context, ports.ListResourceUsageRequest) (domain.ListResourceUsageResponse, error) {
	return domain.ListResourceUsageResponse{Usages: []domain.ResourceUsageResponse{{ResourceUsageID: "usage-1", BusinessID: "business-1", BranchID: "branch-1", ServiceOrderID: "order-1", ResourceID: "resource-1", Quantity: 3}}}, nil
}

func TestManagerCanCreateResourceInAssignedBranch(t *testing.T) {
	app := NewRouter(authServiceStub{role: domain.RoleManager, businessID: "business-1", assignedBranchIDs: []string{"branch-1"}}, branchClientStub{}, resourceClientStub{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/resources", bytes.NewReader([]byte(`{"branch_id":"branch-1","name":"Filter","unit":"pcs","type":"stock"}`)))
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

func TestManagerCannotCreateResourceInUnassignedBranch(t *testing.T) {
	app := NewRouter(authServiceStub{role: domain.RoleManager, businessID: "business-1", assignedBranchIDs: []string{"branch-1"}}, branchClientStub{}, resourceClientStub{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/resources", bytes.NewReader([]byte(`{"branch_id":"branch-2","name":"Filter","unit":"pcs","type":"stock"}`)))
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
