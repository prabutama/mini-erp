package grpcclient

import (
	"context"

	"github.com/isapr/mini-erp/services/api-gateway/internal/domain"
	"github.com/isapr/mini-erp/services/api-gateway/internal/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding"
)

type IdentityClient struct {
	conn *grpc.ClientConn
}

func (c *IdentityClient) CreateUser(ctx context.Context, req ports.CreateUserRequest) (domain.UserResponse, error) {
	var resp domain.UserResponse
	err := c.conn.Invoke(ctx, "/identity.v1.IdentityService/CreateUser", req, &resp)
	return resp, err
}

func (c *IdentityClient) ListUsers(ctx context.Context, req ports.ListUsersRequest) (domain.ListUsersResponse, error) {
	var resp domain.ListUsersResponse
	err := c.conn.Invoke(ctx, "/identity.v1.IdentityService/ListUsers", req, &resp)
	return resp, err
}

func (c *IdentityClient) GetUser(ctx context.Context, req ports.GetUserRequest) (domain.UserResponse, error) {
	var resp domain.UserResponse
	err := c.conn.Invoke(ctx, "/identity.v1.IdentityService/GetUser", req, &resp)
	return resp, err
}

func (c *IdentityClient) UpdateUser(ctx context.Context, req ports.UpdateUserRequest) (domain.UserResponse, error) {
	var resp domain.UserResponse
	err := c.conn.Invoke(ctx, "/identity.v1.IdentityService/UpdateUser", req, &resp)
	return resp, err
}

func (c *IdentityClient) AssignBusinessRole(ctx context.Context, req ports.AssignBusinessRoleRequest) error {
	var resp struct{}
	return c.conn.Invoke(ctx, "/identity.v1.IdentityService/AssignBusinessRole", req, &resp)
}

func NewIdentityClient(ctx context.Context, addr string) (*IdentityClient, error) {
	encoding.RegisterCodec(jsonCodec{})
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.ForceCodec(jsonCodec{})))
	if err != nil {
		return nil, err
	}
	return &IdentityClient{conn: conn}, nil
}

func (c *IdentityClient) Close() error {
	return c.conn.Close()
}

func (c *IdentityClient) SignupTenantAdmin(ctx context.Context, req ports.SignupTenantAdminRequest) (ports.SignupTenantAdminResponse, error) {
	var resp ports.SignupTenantAdminResponse
	err := c.conn.Invoke(ctx, "/identity.v1.IdentityService/SignupTenantAdmin", req, &resp)
	return resp, err
}

func (c *IdentityClient) Login(ctx context.Context, req ports.LoginRequest) (ports.SignupTenantAdminResponse, error) {
	var resp ports.SignupTenantAdminResponse
	err := c.conn.Invoke(ctx, "/identity.v1.IdentityService/Login", req, &resp)
	return resp, err
}

func (c *IdentityClient) GetUserAccessContext(ctx context.Context, req ports.GetUserAccessContextRequest) (ports.AuthContextResponse, error) {
	var resp ports.AuthContextResponse
	err := c.conn.Invoke(ctx, "/identity.v1.IdentityService/GetUserAccessContext", req, &resp)
	return resp, err
}
