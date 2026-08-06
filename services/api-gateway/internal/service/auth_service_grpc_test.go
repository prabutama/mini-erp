package service_test

import (
	"context"
	"testing"

	"github.com/isapr/mini-erp/services/api-gateway/internal/domain"
	"github.com/isapr/mini-erp/services/api-gateway/internal/ports"
	"github.com/isapr/mini-erp/services/api-gateway/internal/service"
)

type organizationClientStub struct {
	req         ports.CreateBusinessRequest
	branchIDs   []string
	assignedReq ports.ListAssignedBranchesRequest
}

func (s *organizationClientStub) CreateBusiness(_ context.Context, req ports.CreateBusinessRequest) (ports.CreateBusinessResponse, error) {
	s.req = req
	return ports.CreateBusinessResponse{BusinessID: "11111111-1111-1111-1111-111111111111", Name: req.Name, Code: "acme", Status: "active", Plan: "free", Timezone: req.Timezone}, nil
}

func (s *organizationClientStub) ListAssignedBranches(_ context.Context, req ports.ListAssignedBranchesRequest) (ports.ListAssignedBranchesResponse, error) {
	s.assignedReq = req
	return ports.ListAssignedBranchesResponse{BranchIDs: s.branchIDs}, nil
}

type identityClientStub struct {
	req        ports.SignupTenantAdminRequest
	accessResp ports.AuthContextResponse
}

func (s *identityClientStub) SignupTenantAdmin(_ context.Context, req ports.SignupTenantAdminRequest) (ports.SignupTenantAdminResponse, error) {
	s.req = req
	return ports.SignupTenantAdminResponse{UserID: "22222222-2222-2222-2222-222222222222", BusinessID: req.BusinessID, Role: string(domain.RoleBusinessAdmin)}, nil
}

func (s *identityClientStub) Login(_ context.Context, _ ports.LoginRequest) (ports.SignupTenantAdminResponse, error) {
	return ports.SignupTenantAdminResponse{}, nil
}

func (s *identityClientStub) GetUserAccessContext(_ context.Context, _ ports.GetUserAccessContextRequest) (ports.AuthContextResponse, error) {
	if s.accessResp.UserID != "" {
		return s.accessResp, nil
	}
	return ports.AuthContextResponse{}, nil
}

func (s *identityClientStub) CreateUser(_ context.Context, _ ports.CreateUserRequest) (domain.UserResponse, error) {
	return domain.UserResponse{}, nil
}

func (s *identityClientStub) ListUsers(_ context.Context, _ ports.ListUsersRequest) (domain.ListUsersResponse, error) {
	return domain.ListUsersResponse{}, nil
}

func (s *identityClientStub) GetUser(_ context.Context, _ ports.GetUserRequest) (domain.UserResponse, error) {
	return domain.UserResponse{}, nil
}

func (s *identityClientStub) UpdateUser(_ context.Context, _ ports.UpdateUserRequest) (domain.UserResponse, error) {
	return domain.UserResponse{}, nil
}

func (s *identityClientStub) AssignBusinessRole(_ context.Context, _ ports.AssignBusinessRoleRequest) error {
	return nil
}

func TestSignupWithClientsPassesBusinessID(t *testing.T) {
	organization := &organizationClientStub{}
	identity := &identityClientStub{}
	auth := service.NewAuthServiceWithClients(identity, organization)

	session, err := auth.Signup(context.Background(), domain.SignupRequest{
		BusinessName:     "Acme Services",
		BusinessTimezone: "Asia/Jakarta",
		AdminFullName:    "Owner",
		AdminEmail:       "owner@example.test",
		AdminPassword:    "secret123",
	}, "request-1")
	if err != nil {
		t.Fatal(err)
	}

	if identity.req.BusinessID == "" {
		t.Fatal("expected business_id passed to identity client")
	}
	if session.User.BusinessID != identity.req.BusinessID {
		t.Fatalf("expected session business_id %q, got %q", identity.req.BusinessID, session.User.BusinessID)
	}
}

func TestGetMeEnrichesManagerAssignedBranches(t *testing.T) {
	organization := &organizationClientStub{branchIDs: []string{"branch-1", "branch-2"}}
	identity := &identityClientStub{accessResp: ports.AuthContextResponse{UserID: "user-1", BusinessID: "business-1", Role: string(domain.RoleManager), RequestID: "request-1"}}
	auth := service.NewAuthServiceWithClients(identity, organization)

	authContext, err := auth.GetMe(context.Background(), "token", "request-1")
	if err != nil {
		t.Fatal(err)
	}

	if organization.assignedReq.UserID != "user-1" || organization.assignedReq.BusinessID != "business-1" {
		t.Fatalf("expected assigned branch lookup for user/business, got %#v", organization.assignedReq)
	}
	if len(authContext.AssignedBranchIDs) != 2 || authContext.AssignedBranchIDs[0] != "branch-1" || authContext.AssignedBranchIDs[1] != "branch-2" {
		t.Fatalf("expected assigned branches from organization, got %#v", authContext.AssignedBranchIDs)
	}
}
