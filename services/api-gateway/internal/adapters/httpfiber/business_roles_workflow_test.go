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

type organizationClientHTTPStub struct{ branchClientStub }

func (organizationClientHTTPStub) CreateBusiness(_ context.Context, req ports.CreateBusinessRequest) (ports.CreateBusinessResponse, error) {
	return ports.CreateBusinessResponse{BusinessID: "business-1", Name: req.Name, Code: "acme", Status: "active", Plan: "free", Timezone: req.Timezone}, nil
}

func (organizationClientHTTPStub) GetBusiness(_ context.Context, req ports.GetBusinessRequest) (domain.BusinessResponse, error) {
	return domain.BusinessResponse{BusinessID: req.BusinessID, Name: "Acme", Code: "acme", Status: "active", Plan: "free", Timezone: "Asia/Jakarta"}, nil
}

func (organizationClientHTTPStub) UpdateBusiness(_ context.Context, req ports.UpdateBusinessRequest) (domain.BusinessResponse, error) {
	return domain.BusinessResponse{BusinessID: req.BusinessID, Name: req.Name, Code: "acme", Status: "active", Plan: "free", Timezone: req.Timezone}, nil
}

func (organizationClientHTTPStub) ListPlatformBusinesses(context.Context, ports.ListPlatformBusinessesRequest) (domain.ListBusinessesResponse, error) {
	return domain.ListBusinessesResponse{Businesses: []domain.BusinessResponse{{BusinessID: "business-1", Name: "Acme", Code: "acme", Status: "active", Plan: "free", Timezone: "Asia/Jakarta"}}}, nil
}

func (organizationClientHTTPStub) UpdatePlatformBusiness(_ context.Context, req ports.UpdatePlatformBusinessRequest) (domain.BusinessResponse, error) {
	return domain.BusinessResponse{BusinessID: req.BusinessID, Name: "Acme", Code: "acme", Status: req.Status, Plan: req.Plan, PlatformNotes: req.PlatformNotes, Timezone: "Asia/Jakarta"}, nil
}

func TestBusinessAdminCanReadAndUpdateCurrentBusiness(t *testing.T) {
	app := NewRouter(authServiceStub{role: domain.RoleBusinessAdmin, businessID: "business-1"}, organizationClientHTTPStub{})
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/businesses/current", nil)
	getReq.Header.Set("Authorization", "Bearer token")
	getResp, err := app.Test(getReq)
	if err != nil {
		t.Fatal(err)
	}
	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", getResp.StatusCode)
	}

	patchReq := httptest.NewRequest(http.MethodPatch, "/api/v1/businesses/current", bytes.NewReader([]byte(`{"name":"Acme Updated","timezone":"Asia/Jakarta"}`)))
	patchReq.Header.Set("Content-Type", "application/json")
	patchReq.Header.Set("Authorization", "Bearer token")
	patchResp, err := app.Test(patchReq)
	if err != nil {
		t.Fatal(err)
	}
	if patchResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", patchResp.StatusCode)
	}
}

func TestPlatformAdminCanUsePlatformBusinessRoutes(t *testing.T) {
	app := NewRouter(authServiceStub{role: domain.RolePlatformAdmin}, organizationClientHTTPStub{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/platform/businesses", nil)
	req.Header.Set("Authorization", "Bearer token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	patchReq := httptest.NewRequest(http.MethodPatch, "/api/v1/platform/businesses/business-1", bytes.NewReader([]byte(`{"status":"suspended","plan":"free","platform_notes":"late payment"}`)))
	patchReq.Header.Set("Content-Type", "application/json")
	patchReq.Header.Set("Authorization", "Bearer token")
	patchResp, err := app.Test(patchReq)
	if err != nil {
		t.Fatal(err)
	}
	if patchResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", patchResp.StatusCode)
	}
}

func TestRolesEndpointReturnsFixedCatalog(t *testing.T) {
	app := NewRouter(authServiceStub{role: domain.RoleBusinessAdmin, businessID: "business-1"}, organizationClientHTTPStub{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/roles", nil)
	req.Header.Set("Authorization", "Bearer token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestWorkflowRoutes(t *testing.T) {
	app := NewRouter(authServiceStub{role: domain.RoleBusinessAdmin, businessID: "business-1"}, organizationClientHTTPStub{}, operationsClientStub{})
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/workflows", bytes.NewReader([]byte(`{"name":"Default","description":"Default workflow"}`)))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer token")
	createResp, err := app.Test(createReq)
	if err != nil {
		t.Fatal(err)
	}
	if createResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", createResp.StatusCode)
	}

	statusReq := httptest.NewRequest(http.MethodPost, "/api/v1/workflows/workflow-1/statuses", bytes.NewReader([]byte(`{"code":"open","name":"Open","is_initial":true}`)))
	statusReq.Header.Set("Content-Type", "application/json")
	statusReq.Header.Set("Authorization", "Bearer token")
	statusResp, err := app.Test(statusReq)
	if err != nil {
		t.Fatal(err)
	}
	if statusResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", statusResp.StatusCode)
	}

	transitionReq := httptest.NewRequest(http.MethodPost, "/api/v1/workflows/workflow-1/transitions", bytes.NewReader([]byte(`{"from_status_code":"open","to_status_code":"in_progress"}`)))
	transitionReq.Header.Set("Content-Type", "application/json")
	transitionReq.Header.Set("Authorization", "Bearer token")
	transitionResp, err := app.Test(transitionReq)
	if err != nil {
		t.Fatal(err)
	}
	if transitionResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", transitionResp.StatusCode)
	}
}
