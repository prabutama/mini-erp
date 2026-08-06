package service_test

import (
	"context"
	"testing"

	"github.com/isapr/mini-erp/services/api-gateway/internal/domain"
	"github.com/isapr/mini-erp/services/api-gateway/internal/ports"
	"github.com/isapr/mini-erp/services/api-gateway/internal/service"
)

type organizationClientStub struct {
	req ports.CreateBusinessRequest
}

func (s *organizationClientStub) CreateBusiness(_ context.Context, req ports.CreateBusinessRequest) (ports.CreateBusinessResponse, error) {
	s.req = req
	return ports.CreateBusinessResponse{BusinessID: "11111111-1111-1111-1111-111111111111", Name: req.Name, Code: "acme", Status: "active", Plan: "free", Timezone: req.Timezone}, nil
}

type identityClientStub struct {
	req ports.SignupTenantAdminRequest
}

func (s *identityClientStub) SignupTenantAdmin(_ context.Context, req ports.SignupTenantAdminRequest) (ports.SignupTenantAdminResponse, error) {
	s.req = req
	return ports.SignupTenantAdminResponse{UserID: "22222222-2222-2222-2222-222222222222", BusinessID: req.BusinessID, Role: string(domain.RoleBusinessAdmin)}, nil
}

func (s *identityClientStub) Login(_ context.Context, _ ports.LoginRequest) (ports.SignupTenantAdminResponse, error) {
	return ports.SignupTenantAdminResponse{}, nil
}

func (s *identityClientStub) GetUserAccessContext(_ context.Context, _ ports.GetUserAccessContextRequest) (ports.AuthContextResponse, error) {
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
