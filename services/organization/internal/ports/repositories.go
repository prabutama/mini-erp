package ports

import (
	"context"

	"github.com/isapr/mini-erp/services/organization/internal/domain"
)

type BusinessRepository interface {
	Create(ctx context.Context, business domain.Business) error
	FindByCode(ctx context.Context, code string) (domain.Business, error)
}

type BranchRepository interface {
	Create(ctx context.Context, branch domain.Branch) error
	FindByID(ctx context.Context, branchID string) (domain.Branch, error)
	ListByBusiness(ctx context.Context, businessID string) ([]domain.Branch, error)
	Update(ctx context.Context, branch domain.Branch) error
}

type PlacementRepository interface {
	Create(ctx context.Context, placement domain.EmployeePlacement) error
}
