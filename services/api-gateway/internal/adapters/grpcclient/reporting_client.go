package grpcclient

import (
	"context"

	"github.com/isapr/mini-erp/services/api-gateway/internal/domain"
	"github.com/isapr/mini-erp/services/api-gateway/internal/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding"
)

type ReportingClient struct{ conn *grpc.ClientConn }

func NewReportingClient(ctx context.Context, addr string) (*ReportingClient, error) {
	encoding.RegisterCodec(jsonCodec{})
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.ForceCodec(jsonCodec{})))
	if err != nil {
		return nil, err
	}
	return &ReportingClient{conn: conn}, nil
}
func (c *ReportingClient) Close() error { return c.conn.Close() }
func (c *ReportingClient) GetAuditEvents(ctx context.Context, req ports.GetAuditEventsRequest) (domain.ListAuditEventsResponse, error) {
	var resp domain.ListAuditEventsResponse
	err := c.conn.Invoke(ctx, "/reporting.v1.ReportingService/GetAuditEvents", req, &resp)
	return resp, err
}
func (c *ReportingClient) GetOperationsSummary(ctx context.Context, req ports.GetOperationsSummaryRequest) (domain.OperationsSummaryReportResponse, error) {
	var resp domain.OperationsSummaryReportResponse
	err := c.conn.Invoke(ctx, "/reporting.v1.ReportingService/GetOperationsSummary", req, &resp)
	return resp, err
}
