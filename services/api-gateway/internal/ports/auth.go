package ports

import (
	"context"

	"github.com/isapr/mini-erp/services/api-gateway/internal/domain"
)

type AuthService interface {
	Signup(ctx context.Context, req domain.SignupRequest, requestID string) (domain.AuthSession, error)
	Login(ctx context.Context, req domain.LoginRequest, requestID string) (domain.AuthSession, error)
	Refresh(ctx context.Context, refreshToken string, requestID string) (domain.AuthSession, error)
	GetMe(ctx context.Context, accessToken string, requestID string) (domain.AuthContext, error)
	Logout(ctx context.Context, accessToken string) error
}

type OrganizationClient interface {
	CreateBusiness(ctx context.Context, req CreateBusinessRequest) (CreateBusinessResponse, error)
}

type IdentityClient interface {
	SignupTenantAdmin(ctx context.Context, req SignupTenantAdminRequest) (SignupTenantAdminResponse, error)
	Login(ctx context.Context, req LoginRequest) (SignupTenantAdminResponse, error)
	GetUserAccessContext(ctx context.Context, req GetUserAccessContextRequest) (AuthContextResponse, error)
	CreateUser(ctx context.Context, req CreateUserRequest) (domain.UserResponse, error)
	ListUsers(ctx context.Context, req ListUsersRequest) (domain.ListUsersResponse, error)
	GetUser(ctx context.Context, req GetUserRequest) (domain.UserResponse, error)
	UpdateUser(ctx context.Context, req UpdateUserRequest) (domain.UserResponse, error)
	AssignBusinessRole(ctx context.Context, req AssignBusinessRoleRequest) error
}

type CreateBusinessRequest struct {
	Name     string `json:"name"`
	Timezone string `json:"timezone"`
}

type CreateBusinessResponse struct {
	BusinessID string `json:"business_id"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	Status     string `json:"status"`
	Plan       string `json:"plan"`
	Timezone   string `json:"timezone"`
}

type SignupTenantAdminRequest struct {
	BusinessID    string `json:"business_id"`
	AdminEmail    string `json:"admin_email"`
	AdminPassword string `json:"admin_password"`
	AdminFullName string `json:"admin_full_name"`
	RequestID     string `json:"request_id"`
}

type SignupTenantAdminResponse struct {
	AccessToken       string   `json:"access_token"`
	RefreshToken      string   `json:"refresh_token"`
	UserID            string   `json:"user_id"`
	BusinessID        string   `json:"business_id"`
	Role              string   `json:"role"`
	Permissions       []string `json:"permissions"`
	AssignedBranchIDs []string `json:"assigned_branch_ids"`
	RequestID         string   `json:"request_id"`
}

type LoginRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	RequestID string `json:"request_id"`
}

type GetUserAccessContextRequest struct {
	AccessToken string `json:"access_token"`
	RequestID   string `json:"request_id"`
}

type AuthContextResponse struct {
	UserID            string   `json:"user_id"`
	BusinessID        string   `json:"business_id,omitempty"`
	Role              string   `json:"role"`
	Permissions       []string `json:"permissions"`
	AssignedBranchIDs []string `json:"assigned_branch_ids,omitempty"`
	RequestID         string   `json:"request_id"`
}

type CreateUserRequest struct {
	BusinessID string `json:"business_id"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	FullName   string `json:"full_name"`
	Role       string `json:"role"`
}

type ListUsersRequest struct {
	BusinessID string `json:"business_id"`
}

type GetUserRequest struct {
	UserID string `json:"user_id"`
}

type UpdateUserRequest struct {
	UserID   string `json:"user_id"`
	FullName string `json:"full_name"`
	Status   string `json:"status"`
}

type AssignBusinessRoleRequest struct {
	UserID     string `json:"user_id"`
	BusinessID string `json:"business_id"`
	Role       string `json:"role"`
}
