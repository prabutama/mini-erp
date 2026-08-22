package grpcclient

import (
	"context"

	"github.com/isapr/mini-erp/services/api-gateway/internal/domain"
	"github.com/isapr/mini-erp/services/api-gateway/internal/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding"
)

type OperationsClient struct{ conn *grpc.ClientConn }

func NewOperationsClient(ctx context.Context, addr string) (*OperationsClient, error) {
	encoding.RegisterCodec(jsonCodec{})
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.ForceCodec(jsonCodec{})))
	if err != nil {
		return nil, err
	}
	return &OperationsClient{conn: conn}, nil
}

func (c *OperationsClient) Close() error { return c.conn.Close() }

func (c *OperationsClient) CreateServiceDefinition(ctx context.Context, req ports.CreateServiceDefinitionRequest) (domain.ServiceDefinitionResponse, error) {
	var resp domain.ServiceDefinitionResponse
	err := c.conn.Invoke(ctx, "/operations.v1.OperationsService/CreateServiceDefinition", req, &resp)
	return resp, err
}

func (c *OperationsClient) ListServiceDefinitions(ctx context.Context, req ports.ListServiceDefinitionsRequest) (domain.ListServiceDefinitionsResponse, error) {
	var resp domain.ListServiceDefinitionsResponse
	err := c.conn.Invoke(ctx, "/operations.v1.OperationsService/ListServiceDefinitions", req, &resp)
	return resp, err
}

func (c *OperationsClient) CreateServiceOrder(ctx context.Context, req ports.CreateServiceOrderRequest) (domain.ServiceOrderResponse, error) {
	var resp domain.ServiceOrderResponse
	err := c.conn.Invoke(ctx, "/operations.v1.OperationsService/CreateServiceOrder", req, &resp)
	return resp, err
}

func (c *OperationsClient) ListServiceOrders(ctx context.Context, req ports.ListServiceOrdersRequest) (domain.ListServiceOrdersResponse, error) {
	var resp domain.ListServiceOrdersResponse
	err := c.conn.Invoke(ctx, "/operations.v1.OperationsService/ListServiceOrders", req, &resp)
	return resp, err
}

func (c *OperationsClient) ServiceOrderSummary(ctx context.Context, req ports.ServiceOrderSummaryRequest) (domain.ServiceOrderSummaryResponse, error) {
	var resp domain.ServiceOrderSummaryResponse
	err := c.conn.Invoke(ctx, "/operations.v1.OperationsService/ServiceOrderSummary", req, &resp)
	return resp, err
}

func (c *OperationsClient) GetServiceOrder(ctx context.Context, req ports.GetServiceOrderRequest) (domain.ServiceOrderResponse, error) {
	var resp domain.ServiceOrderResponse
	err := c.conn.Invoke(ctx, "/operations.v1.OperationsService/GetServiceOrder", req, &resp)
	return resp, err
}

func (c *OperationsClient) TransitionServiceOrder(ctx context.Context, req ports.TransitionServiceOrderRequest) (domain.ServiceOrderResponse, error) {
	var resp domain.ServiceOrderResponse
	err := c.conn.Invoke(ctx, "/operations.v1.OperationsService/TransitionServiceOrder", req, &resp)
	return resp, err
}

func (c *OperationsClient) AssignServiceOrder(ctx context.Context, req ports.AssignServiceOrderRequest) (domain.ServiceOrderAssignmentResponse, error) {
	var resp domain.ServiceOrderAssignmentResponse
	err := c.conn.Invoke(ctx, "/operations.v1.OperationsService/AssignServiceOrder", req, &resp)
	return resp, err
}

func (c *OperationsClient) ListAssignedServiceOrders(ctx context.Context, req ports.ListAssignedServiceOrdersRequest) (domain.ListServiceOrdersResponse, error) {
	var resp domain.ListServiceOrdersResponse
	err := c.conn.Invoke(ctx, "/operations.v1.OperationsService/ListAssignedServiceOrders", req, &resp)
	return resp, err
}

func (c *OperationsClient) ListServiceOrderAssignments(ctx context.Context, req ports.ListServiceOrderAssignmentsRequest) (domain.ListServiceOrderAssignmentsResponse, error) {
	var resp domain.ListServiceOrderAssignmentsResponse
	err := c.conn.Invoke(ctx, "/operations.v1.OperationsService/ListServiceOrderAssignments", req, &resp)
	return resp, err
}

func (c *OperationsClient) CreateWorkflow(ctx context.Context, req ports.CreateWorkflowRequest) (domain.WorkflowResponse, error) {
	var resp domain.WorkflowResponse
	err := c.conn.Invoke(ctx, "/operations.v1.OperationsService/CreateWorkflow", req, &resp)
	return resp, err
}
func (c *OperationsClient) ListWorkflows(ctx context.Context, req ports.ListWorkflowsRequest) (domain.ListWorkflowsResponse, error) {
	var resp domain.ListWorkflowsResponse
	err := c.conn.Invoke(ctx, "/operations.v1.OperationsService/ListWorkflows", req, &resp)
	return resp, err
}
func (c *OperationsClient) GetWorkflow(ctx context.Context, req ports.GetWorkflowRequest) (domain.WorkflowResponse, error) {
	var resp domain.WorkflowResponse
	err := c.conn.Invoke(ctx, "/operations.v1.OperationsService/GetWorkflow", req, &resp)
	return resp, err
}
func (c *OperationsClient) UpdateWorkflow(ctx context.Context, req ports.UpdateWorkflowRequest) (domain.WorkflowResponse, error) {
	var resp domain.WorkflowResponse
	err := c.conn.Invoke(ctx, "/operations.v1.OperationsService/UpdateWorkflow", req, &resp)
	return resp, err
}
func (c *OperationsClient) CreateWorkflowStatus(ctx context.Context, req ports.CreateWorkflowStatusRequest) (domain.WorkflowStatusResponse, error) {
	var resp domain.WorkflowStatusResponse
	err := c.conn.Invoke(ctx, "/operations.v1.OperationsService/CreateWorkflowStatus", req, &resp)
	return resp, err
}
func (c *OperationsClient) CreateWorkflowTransition(ctx context.Context, req ports.CreateWorkflowTransitionRequest) (domain.WorkflowTransitionResponse, error) {
	var resp domain.WorkflowTransitionResponse
	err := c.conn.Invoke(ctx, "/operations.v1.OperationsService/CreateWorkflowTransition", req, &resp)
	return resp, err
}
