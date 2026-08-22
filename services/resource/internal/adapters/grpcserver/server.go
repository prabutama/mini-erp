package grpcserver

import (
	"context"
	"errors"

	"github.com/isapr/mini-erp/services/resource/internal/application"
	"github.com/isapr/mini-erp/services/resource/internal/domain"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/status"
)

const serviceName = "resource.v1.ResourceService"

type CreateResourceRequest struct {
	BusinessID string `json:"business_id"`
	BranchID   string `json:"branch_id"`
	Name       string `json:"name"`
	Unit       string `json:"unit"`
	Type       string `json:"type"`
}

type ResourceResponse struct {
	ResourceID string `json:"resource_id"`
	BusinessID string `json:"business_id"`
	BranchID   string `json:"branch_id"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	Unit       string `json:"unit"`
	Type       string `json:"type"`
	Status     string `json:"status"`
}

type ListResourcesRequest struct {
	BusinessID string `json:"business_id"`
	BranchID   string `json:"branch_id"`
}

type ListResourcesResponse struct {
	Resources []ResourceResponse `json:"resources"`
}
type GetResourceRequest struct {
	ResourceID string `json:"resource_id"`
}

type RecordStockMovementRequest struct {
	BusinessID     string  `json:"business_id"`
	BranchID       string  `json:"branch_id"`
	ResourceID     string  `json:"resource_id"`
	MovementType   string  `json:"movement_type"`
	Reason         string  `json:"reason"`
	ServiceOrderID string  `json:"service_order_id"`
	ActorUserID    string  `json:"actor_user_id"`
	RequestID      string  `json:"request_id"`
	Quantity       float64 `json:"quantity"`
}

type StockMovementResponse struct {
	StockMovementID string  `json:"stock_movement_id"`
	BusinessID      string  `json:"business_id"`
	BranchID        string  `json:"branch_id"`
	ResourceID      string  `json:"resource_id"`
	MovementType    string  `json:"movement_type"`
	Reason          string  `json:"reason"`
	ServiceOrderID  string  `json:"service_order_id"`
	ActorUserID     string  `json:"actor_user_id"`
	Quantity        float64 `json:"quantity"`
}
type GetResourceAvailabilityRequest struct {
	ResourceID string `json:"resource_id"`
}
type ResourceAvailabilityResponse struct {
	ResourceID string  `json:"resource_id"`
	Quantity   float64 `json:"quantity"`
}

type RecordResourceUsageRequest struct {
	BusinessID       string  `json:"business_id"`
	BranchID         string  `json:"branch_id"`
	ServiceOrderID   string  `json:"service_order_id"`
	ResourceID       string  `json:"resource_id"`
	Quantity         float64 `json:"quantity"`
	Reason           string  `json:"reason"`
	RecordedByUserID string  `json:"recorded_by_user_id"`
	RequestID        string  `json:"request_id"`
}

type ResourceUsageResponse struct {
	ResourceUsageID  string  `json:"resource_usage_id"`
	BusinessID       string  `json:"business_id"`
	BranchID         string  `json:"branch_id"`
	ServiceOrderID   string  `json:"service_order_id"`
	ResourceID       string  `json:"resource_id"`
	Quantity         float64 `json:"quantity"`
	Reason           string  `json:"reason"`
	RecordedByUserID string  `json:"recorded_by_user_id"`
	StockMovementID  string  `json:"stock_movement_id"`
}

type ListResourceUsageRequest struct {
	ServiceOrderID string `json:"service_order_id"`
}
type ListResourceUsageResponse struct {
	Usages []ResourceUsageResponse `json:"usages"`
}

type Server struct{ resources *application.ResourceService }

type resourceServiceServer interface {
	CreateResource(context.Context, CreateResourceRequest) (ResourceResponse, error)
	ListResources(context.Context, ListResourcesRequest) (ListResourcesResponse, error)
	GetResource(context.Context, GetResourceRequest) (ResourceResponse, error)
	RecordStockMovement(context.Context, RecordStockMovementRequest) (StockMovementResponse, error)
	GetResourceAvailability(context.Context, GetResourceAvailabilityRequest) (ResourceAvailabilityResponse, error)
	RecordResourceUsage(context.Context, RecordResourceUsageRequest) (ResourceUsageResponse, error)
	ListResourceUsage(context.Context, ListResourceUsageRequest) (ListResourceUsageResponse, error)
}

func New(resources *application.ResourceService) *Server {
	encoding.RegisterCodec(jsonCodec{})
	return &Server{resources: resources}
}

func (s *Server) Register(grpcServer *grpc.Server) {
	grpcServer.RegisterService(&grpc.ServiceDesc{ServiceName: serviceName, HandlerType: (*resourceServiceServer)(nil), Methods: []grpc.MethodDesc{{MethodName: "CreateResource", Handler: createResourceHandler}, {MethodName: "ListResources", Handler: listResourcesHandler}, {MethodName: "GetResource", Handler: getResourceHandler}, {MethodName: "RecordStockMovement", Handler: recordStockMovementHandler}, {MethodName: "GetResourceAvailability", Handler: getResourceAvailabilityHandler}, {MethodName: "RecordResourceUsage", Handler: recordResourceUsageHandler}, {MethodName: "ListResourceUsage", Handler: listResourceUsageHandler}}}, s)
}

func createResourceHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req CreateResourceRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(resourceServiceServer).CreateResource(ctx, req)
}
func listResourcesHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req ListResourcesRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(resourceServiceServer).ListResources(ctx, req)
}
func getResourceHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req GetResourceRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(resourceServiceServer).GetResource(ctx, req)
}
func recordStockMovementHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req RecordStockMovementRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(resourceServiceServer).RecordStockMovement(ctx, req)
}
func getResourceAvailabilityHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req GetResourceAvailabilityRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(resourceServiceServer).GetResourceAvailability(ctx, req)
}

func recordResourceUsageHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req RecordResourceUsageRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(resourceServiceServer).RecordResourceUsage(ctx, req)
}

func listResourceUsageHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req ListResourceUsageRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(resourceServiceServer).ListResourceUsage(ctx, req)
}

func (s *Server) CreateResource(ctx context.Context, req CreateResourceRequest) (ResourceResponse, error) {
	resource, err := s.resources.Create(ctx, application.CreateResourceInput{BusinessID: req.BusinessID, BranchID: req.BranchID, Name: req.Name, Unit: req.Unit, Type: req.Type})
	if err != nil {
		return ResourceResponse{}, mapError(err)
	}
	return resourceResponse(resource), nil
}
func (s *Server) ListResources(ctx context.Context, req ListResourcesRequest) (ListResourcesResponse, error) {
	resources, err := s.resources.List(ctx, req.BusinessID, req.BranchID)
	if err != nil {
		return ListResourcesResponse{}, mapError(err)
	}
	resp := ListResourcesResponse{Resources: []ResourceResponse{}}
	for _, resource := range resources {
		resp.Resources = append(resp.Resources, resourceResponse(resource))
	}
	return resp, nil
}
func (s *Server) GetResource(ctx context.Context, req GetResourceRequest) (ResourceResponse, error) {
	resource, err := s.resources.Get(ctx, req.ResourceID)
	if err != nil {
		return ResourceResponse{}, mapError(err)
	}
	return resourceResponse(resource), nil
}
func (s *Server) RecordStockMovement(ctx context.Context, req RecordStockMovementRequest) (StockMovementResponse, error) {
	movement, err := s.resources.RecordStockMovement(ctx, application.RecordStockMovementInput{BusinessID: req.BusinessID, BranchID: req.BranchID, ResourceID: req.ResourceID, MovementType: req.MovementType, Quantity: req.Quantity, Reason: req.Reason, ServiceOrderID: req.ServiceOrderID, ActorUserID: req.ActorUserID, RequestID: req.RequestID})
	if err != nil {
		return StockMovementResponse{}, mapError(err)
	}
	return movementResponse(movement), nil
}
func (s *Server) GetResourceAvailability(ctx context.Context, req GetResourceAvailabilityRequest) (ResourceAvailabilityResponse, error) {
	availability, err := s.resources.Availability(ctx, req.ResourceID)
	if err != nil {
		return ResourceAvailabilityResponse{}, mapError(err)
	}
	return ResourceAvailabilityResponse{ResourceID: availability.ResourceID.String(), Quantity: availability.Quantity}, nil
}

func (s *Server) RecordResourceUsage(ctx context.Context, req RecordResourceUsageRequest) (ResourceUsageResponse, error) {
	usage, err := s.resources.RecordResourceUsage(ctx, application.RecordResourceUsageInput{BusinessID: req.BusinessID, BranchID: req.BranchID, ServiceOrderID: req.ServiceOrderID, ResourceID: req.ResourceID, Quantity: req.Quantity, Reason: req.Reason, RecordedByUserID: req.RecordedByUserID, RequestID: req.RequestID})
	if err != nil {
		return ResourceUsageResponse{}, mapError(err)
	}
	return usageResponse(usage), nil
}

func (s *Server) ListResourceUsage(ctx context.Context, req ListResourceUsageRequest) (ListResourceUsageResponse, error) {
	usages, err := s.resources.ListUsageByServiceOrder(ctx, req.ServiceOrderID)
	if err != nil {
		return ListResourceUsageResponse{}, mapError(err)
	}
	response := ListResourceUsageResponse{Usages: []ResourceUsageResponse{}}
	for _, usage := range usages {
		response.Usages = append(response.Usages, usageResponse(usage))
	}
	return response, nil
}

func mapError(err error) error {
	if errors.Is(err, application.ErrValidation) {
		return status.Error(codes.InvalidArgument, "validation failed")
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return status.Error(codes.NotFound, "not found")
	}
	return status.Error(codes.Internal, err.Error())
}
func resourceResponse(resource domain.Resource) ResourceResponse {
	return ResourceResponse{ResourceID: resource.ID.String(), BusinessID: resource.BusinessID.String(), BranchID: resource.BranchID.String(), Name: resource.Name, Code: resource.Code, Unit: resource.Unit, Type: resource.Type, Status: resource.Status}
}
func movementResponse(movement domain.StockMovement) StockMovementResponse {
	return StockMovementResponse{StockMovementID: movement.ID.String(), BusinessID: movement.BusinessID.String(), BranchID: movement.BranchID.String(), ResourceID: movement.ResourceID.String(), MovementType: movement.MovementType, Quantity: movement.Quantity, Reason: movement.Reason, ServiceOrderID: movement.ServiceOrderID.String(), ActorUserID: movement.ActorUserID.String()}
}

func usageResponse(usage domain.ResourceUsage) ResourceUsageResponse {
	return ResourceUsageResponse{ResourceUsageID: usage.ID.String(), BusinessID: usage.BusinessID.String(), BranchID: usage.BranchID.String(), ServiceOrderID: usage.ServiceOrderID.String(), ResourceID: usage.ResourceID.String(), Quantity: usage.Quantity, Reason: usage.Reason, RecordedByUserID: usage.RecordedByUserID.String(), StockMovementID: usage.StockMovementID.String()}
}
