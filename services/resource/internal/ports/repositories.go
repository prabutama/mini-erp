package ports

import (
	"context"

	"github.com/isapr/mini-erp/services/resource/internal/domain"
)

type ResourceRepository interface {
	Create(ctx context.Context, resource domain.Resource) error
	List(ctx context.Context, businessID string, branchID string) ([]domain.Resource, error)
	FindByID(ctx context.Context, resourceID string) (domain.Resource, error)
	RecordStockMovement(ctx context.Context, movement domain.StockMovement) error
	Availability(ctx context.Context, resourceID string) (domain.ResourceAvailability, error)
	RecordResourceUsage(ctx context.Context, usage domain.ResourceUsage, movement domain.StockMovement) error
	ListUsageByServiceOrder(ctx context.Context, serviceOrderID string) ([]domain.ResourceUsage, error)
}
