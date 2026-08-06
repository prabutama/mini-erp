package httpfiber

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/isapr/mini-erp/services/api-gateway/internal/domain"
	"github.com/isapr/mini-erp/services/api-gateway/internal/ports"
	"github.com/isapr/mini-erp/services/api-gateway/internal/service"
)

type branchClientStub struct{}

func (branchClientStub) CreateBranch(_ context.Context, req ports.CreateBranchRequest) (domain.BranchResponse, error) {
	return domain.BranchResponse{BranchID: "branch-1", BusinessID: req.BusinessID, Name: req.Name, Code: "main", Status: "active"}, nil
}

func (branchClientStub) GetBranch(_ context.Context, req ports.GetBranchRequest) (domain.BranchResponse, error) {
	return domain.BranchResponse{BranchID: req.BranchID, BusinessID: "business-1", Name: "Main", Code: "main", Status: "active"}, nil
}

func (branchClientStub) ListBranches(_ context.Context, _ ports.ListBranchesRequest) (domain.ListBranchesResponse, error) {
	return domain.ListBranchesResponse{Branches: []domain.BranchResponse{{BranchID: "branch-1", BusinessID: "business-1", Name: "Main", Code: "main", Status: "active"}, {BranchID: "branch-2", BusinessID: "business-1", Name: "Second", Code: "second", Status: "active"}}}, nil
}

func (branchClientStub) UpdateBranch(_ context.Context, req ports.UpdateBranchRequest) (domain.BranchResponse, error) {
	return domain.BranchResponse{BranchID: req.BranchID, BusinessID: "business-1", Name: req.Name, Code: "main", Status: req.Status}, nil
}

func (branchClientStub) CreateEmployeePlacement(_ context.Context, req ports.CreateEmployeePlacementRequest) (domain.PlacementResponse, error) {
	return domain.PlacementResponse{PlacementID: "placement-1", UserID: req.UserID, BusinessID: req.BusinessID, BranchID: req.BranchID, Position: req.Position, EmploymentType: req.EmploymentType, Status: "active"}, nil
}

func (branchClientStub) ListAssignedBranches(_ context.Context, req ports.ListAssignedBranchesRequest) (ports.ListAssignedBranchesResponse, error) {
	return ports.ListAssignedBranchesResponse{BranchIDs: []string{"branch-1"}, Placements: []domain.PlacementResponse{{PlacementID: "placement-1", UserID: req.UserID, BusinessID: req.BusinessID, BranchID: "branch-1", Status: "active"}}}, nil
}

func TestCreateBranch(t *testing.T) {
	auth := service.NewAuthService()
	app := NewRouter(auth, branchClientStub{})

	signupBody, _ := json.Marshal(domain.SignupRequest{BusinessName: "Acme", BusinessTimezone: "Asia/Jakarta", AdminFullName: "Owner", AdminEmail: "owner-branch@example.test", AdminPassword: "secret123"})
	signupReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewReader(signupBody))
	signupReq.Header.Set("Content-Type", "application/json")
	signupResp, err := app.Test(signupReq)
	if err != nil {
		t.Fatal(err)
	}
	var session domain.AuthSession
	if err := json.NewDecoder(signupResp.Body).Decode(&session); err != nil {
		t.Fatal(err)
	}

	branchBody, _ := json.Marshal(domain.CreateBranchRequest{Name: "Main Branch"})
	branchReq := httptest.NewRequest(http.MethodPost, "/api/v1/branches", bytes.NewReader(branchBody))
	branchReq.Header.Set("Content-Type", "application/json")
	branchReq.Header.Set("Authorization", "Bearer "+session.AccessToken)
	branchResp, err := app.Test(branchReq)
	if err != nil {
		t.Fatal(err)
	}
	if branchResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", branchResp.StatusCode)
	}
}
