package ports

import (
	"context"

	"github.com/isapr/mini-erp/services/api-gateway/internal/domain"
)

type ResourceClient interface {
	CreateResource(ctx context.Context, req CreateResourceRequest) (domain.ResourceResponse, error)
	ListResources(ctx context.Context, req ListResourcesRequest) (domain.ListResourcesResponse, error)
	GetResource(ctx context.Context, req GetResourceRequest) (domain.ResourceResponse, error)
	RecordStockMovement(ctx context.Context, req RecordStockMovementRequest) (domain.StockMovementResponse, error)
	GetResourceAvailability(ctx context.Context, req GetResourceAvailabilityRequest) (domain.ResourceAvailabilityResponse, error)
	RecordResourceUsage(ctx context.Context, req RecordResourceUsageRequest) (domain.ResourceUsageResponse, error)
	ListResourceUsage(ctx context.Context, req ListResourceUsageRequest) (domain.ListResourceUsageResponse, error)
}

type CreateResourceRequest struct {
	BusinessID string `json:"business_id"`
	BranchID   string `json:"branch_id"`
	Name       string `json:"name"`
	Unit       string `json:"unit"`
	Type       string `json:"type"`
}

type ListResourcesRequest struct {
	BusinessID string `json:"business_id"`
	BranchID   string `json:"branch_id"`
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
type GetResourceAvailabilityRequest struct {
	ResourceID string `json:"resource_id"`
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

type ListResourceUsageRequest struct {
	ServiceOrderID string `json:"service_order_id"`
}
