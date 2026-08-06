package grpcclient

import (
	"context"

	"github.com/isapr/mini-erp/services/api-gateway/internal/domain"
	"github.com/isapr/mini-erp/services/api-gateway/internal/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding"
)

type OrganizationClient struct {
	conn *grpc.ClientConn
}

func (c *OrganizationClient) CreateBranch(ctx context.Context, req ports.CreateBranchRequest) (domain.BranchResponse, error) {
	var resp domain.BranchResponse
	err := c.conn.Invoke(ctx, "/organization.v1.OrganizationService/CreateBranch", req, &resp)
	return resp, err
}

func (c *OrganizationClient) GetBranch(ctx context.Context, req ports.GetBranchRequest) (domain.BranchResponse, error) {
	var resp domain.BranchResponse
	err := c.conn.Invoke(ctx, "/organization.v1.OrganizationService/GetBranch", req, &resp)
	return resp, err
}

func (c *OrganizationClient) ListBranches(ctx context.Context, req ports.ListBranchesRequest) (domain.ListBranchesResponse, error) {
	var resp domain.ListBranchesResponse
	err := c.conn.Invoke(ctx, "/organization.v1.OrganizationService/ListBranches", req, &resp)
	return resp, err
}

func (c *OrganizationClient) UpdateBranch(ctx context.Context, req ports.UpdateBranchRequest) (domain.BranchResponse, error) {
	var resp domain.BranchResponse
	err := c.conn.Invoke(ctx, "/organization.v1.OrganizationService/UpdateBranch", req, &resp)
	return resp, err
}

func (c *OrganizationClient) CreateEmployeePlacement(ctx context.Context, req ports.CreateEmployeePlacementRequest) (domain.PlacementResponse, error) {
	var resp domain.PlacementResponse
	err := c.conn.Invoke(ctx, "/organization.v1.OrganizationService/CreateEmployeePlacement", req, &resp)
	return resp, err
}

func (c *OrganizationClient) ListAssignedBranches(ctx context.Context, req ports.ListAssignedBranchesRequest) (ports.ListAssignedBranchesResponse, error) {
	var resp ports.ListAssignedBranchesResponse
	err := c.conn.Invoke(ctx, "/organization.v1.OrganizationService/ListAssignedBranches", req, &resp)
	return resp, err
}

func NewOrganizationClient(ctx context.Context, addr string) (*OrganizationClient, error) {
	encoding.RegisterCodec(jsonCodec{})
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.ForceCodec(jsonCodec{})))
	if err != nil {
		return nil, err
	}
	return &OrganizationClient{conn: conn}, nil
}

func (c *OrganizationClient) Close() error {
	return c.conn.Close()
}

func (c *OrganizationClient) CreateBusiness(ctx context.Context, req ports.CreateBusinessRequest) (ports.CreateBusinessResponse, error) {
	var resp ports.CreateBusinessResponse
	err := c.conn.Invoke(ctx, "/organization.v1.OrganizationService/CreateBusiness", req, &resp)
	return resp, err
}
