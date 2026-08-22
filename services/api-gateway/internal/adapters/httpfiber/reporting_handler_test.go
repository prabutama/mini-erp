package httpfiber

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/isapr/mini-erp/services/api-gateway/internal/domain"
	"github.com/isapr/mini-erp/services/api-gateway/internal/ports"
)

type reportingClientStub struct{}

func (reportingClientStub) GetAuditEvents(context.Context, ports.GetAuditEventsRequest) (domain.ListAuditEventsResponse, error) {
	return domain.ListAuditEventsResponse{Events: []domain.AuditEventResponse{}}, nil
}
func (reportingClientStub) GetOperationsSummary(context.Context, ports.GetOperationsSummaryRequest) (domain.OperationsSummaryReportResponse, error) {
	return domain.OperationsSummaryReportResponse{BusinessID: "business-1", BranchID: "branch-1", SnapshotDate: "2026-08-22"}, nil
}

func TestBusinessAdminCanReadAuditEvents(t *testing.T) {
	app := NewRouter(authServiceStub{role: domain.RoleBusinessAdmin, businessID: "business-1"}, branchClientStub{}, reportingClientStub{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/reports/audit-events", nil)
	req.Header.Set("Authorization", "Bearer token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestManagerCannotReadUnassignedBranchReport(t *testing.T) {
	app := NewRouter(authServiceStub{role: domain.RoleManager, businessID: "business-1", assignedBranchIDs: []string{"branch-1"}}, branchClientStub{}, reportingClientStub{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/reports/operations-summary?branch_id=branch-2", nil)
	req.Header.Set("Authorization", "Bearer token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", resp.StatusCode)
	}
}

func TestStaffCannotReadReports(t *testing.T) {
	app := NewRouter(authServiceStub{role: domain.RoleStaff, businessID: "business-1", assignedBranchIDs: []string{"branch-1"}}, branchClientStub{}, reportingClientStub{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/reports/audit-events?branch_id=branch-1", nil)
	req.Header.Set("Authorization", "Bearer token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", resp.StatusCode)
	}
}
