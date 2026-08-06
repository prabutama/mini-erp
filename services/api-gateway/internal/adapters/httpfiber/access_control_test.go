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

type authServiceStub struct {
	role       domain.Role
	businessID string
}

func (s authServiceStub) Signup(context.Context, domain.SignupRequest, string) (domain.AuthSession, error) {
	return domain.AuthSession{}, nil
}

func (s authServiceStub) Login(context.Context, domain.LoginRequest, string) (domain.AuthSession, error) {
	return domain.AuthSession{}, nil
}

func (s authServiceStub) Refresh(context.Context, string, string) (domain.AuthSession, error) {
	return domain.AuthSession{}, nil
}

func (s authServiceStub) GetMe(context.Context, string, string) (domain.AuthContext, error) {
	return domain.AuthContext{UserID: "user-1", BusinessID: s.businessID, Role: s.role, Permissions: []string{}, RequestID: "request-1"}, nil
}

func (s authServiceStub) Logout(context.Context, string) error {
	return nil
}

type identityClientStubForHTTP struct{}

func (identityClientStubForHTTP) SignupTenantAdmin(context.Context, ports.SignupTenantAdminRequest) (ports.SignupTenantAdminResponse, error) {
	return ports.SignupTenantAdminResponse{}, nil
}

func (identityClientStubForHTTP) Login(context.Context, ports.LoginRequest) (ports.SignupTenantAdminResponse, error) {
	return ports.SignupTenantAdminResponse{}, nil
}

func (identityClientStubForHTTP) GetUserAccessContext(context.Context, ports.GetUserAccessContextRequest) (ports.AuthContextResponse, error) {
	return ports.AuthContextResponse{}, nil
}

func (identityClientStubForHTTP) CreateUser(context.Context, ports.CreateUserRequest) (domain.UserResponse, error) {
	return domain.UserResponse{}, nil
}

func (identityClientStubForHTTP) ListUsers(context.Context, ports.ListUsersRequest) (domain.ListUsersResponse, error) {
	return domain.ListUsersResponse{}, nil
}

func (identityClientStubForHTTP) GetUser(context.Context, ports.GetUserRequest) (domain.UserResponse, error) {
	return domain.UserResponse{}, nil
}

func (identityClientStubForHTTP) UpdateUser(context.Context, ports.UpdateUserRequest) (domain.UserResponse, error) {
	return domain.UserResponse{}, nil
}

func (identityClientStubForHTTP) AssignBusinessRole(context.Context, ports.AssignBusinessRoleRequest) error {
	return nil
}

func TestManagerCannotCreateBranch(t *testing.T) {
	app := NewRouter(authServiceStub{role: domain.RoleManager, businessID: "business-1"}, branchClientStub{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/branches", bytes.NewReader([]byte(`{"name":"Main"}`)))
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

func TestPlatformAdminCannotCreateBranch(t *testing.T) {
	app := NewRouter(authServiceStub{role: domain.RolePlatformAdmin}, branchClientStub{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/branches", bytes.NewReader([]byte(`{"name":"Main"}`)))
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

func TestBusinessAdminCanCreateBranch(t *testing.T) {
	app := NewRouter(authServiceStub{role: domain.RoleBusinessAdmin, businessID: "business-1"}, branchClientStub{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/branches", bytes.NewReader([]byte(`{"name":"Main"}`)))
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

func TestManagerCannotCreateUser(t *testing.T) {
	app := NewRouter(authServiceStub{role: domain.RoleManager, businessID: "business-1"}, branchClientStub{}, identityClientStubForHTTP{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader([]byte(`{"email":"staff@example.test","password":"secret123","full_name":"Staff","role":"Staff"}`)))
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
