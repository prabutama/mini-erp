package grpcclient

import (
	"context"

	"github.com/isapr/mini-erp/services/api-gateway/internal/domain"
	"github.com/isapr/mini-erp/services/api-gateway/internal/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding"
)

type ResourceClient struct{ conn *grpc.ClientConn }

func NewResourceClient(ctx context.Context, addr string) (*ResourceClient, error) {
	encoding.RegisterCodec(jsonCodec{})
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.ForceCodec(jsonCodec{})))
	if err != nil {
		return nil, err
	}
	return &ResourceClient{conn: conn}, nil
}

func (c *ResourceClient) Close() error { return c.conn.Close() }
func (c *ResourceClient) CreateResource(ctx context.Context, req ports.CreateResourceRequest) (domain.ResourceResponse, error) {
	var resp domain.ResourceResponse
	err := c.conn.Invoke(ctx, "/resource.v1.ResourceService/CreateResource", req, &resp)
	return resp, err
}
func (c *ResourceClient) ListResources(ctx context.Context, req ports.ListResourcesRequest) (domain.ListResourcesResponse, error) {
	var resp domain.ListResourcesResponse
	err := c.conn.Invoke(ctx, "/resource.v1.ResourceService/ListResources", req, &resp)
	return resp, err
}
func (c *ResourceClient) GetResource(ctx context.Context, req ports.GetResourceRequest) (domain.ResourceResponse, error) {
	var resp domain.ResourceResponse
	err := c.conn.Invoke(ctx, "/resource.v1.ResourceService/GetResource", req, &resp)
	return resp, err
}
func (c *ResourceClient) RecordStockMovement(ctx context.Context, req ports.RecordStockMovementRequest) (domain.StockMovementResponse, error) {
	var resp domain.StockMovementResponse
	err := c.conn.Invoke(ctx, "/resource.v1.ResourceService/RecordStockMovement", req, &resp)
	return resp, err
}
func (c *ResourceClient) GetResourceAvailability(ctx context.Context, req ports.GetResourceAvailabilityRequest) (domain.ResourceAvailabilityResponse, error) {
	var resp domain.ResourceAvailabilityResponse
	err := c.conn.Invoke(ctx, "/resource.v1.ResourceService/GetResourceAvailability", req, &resp)
	return resp, err
}

func (c *ResourceClient) RecordResourceUsage(ctx context.Context, req ports.RecordResourceUsageRequest) (domain.ResourceUsageResponse, error) {
	var resp domain.ResourceUsageResponse
	err := c.conn.Invoke(ctx, "/resource.v1.ResourceService/RecordResourceUsage", req, &resp)
	return resp, err
}

func (c *ResourceClient) ListResourceUsage(ctx context.Context, req ports.ListResourceUsageRequest) (domain.ListResourceUsageResponse, error) {
	var resp domain.ListResourceUsageResponse
	err := c.conn.Invoke(ctx, "/resource.v1.ResourceService/ListResourceUsage", req, &resp)
	return resp, err
}
