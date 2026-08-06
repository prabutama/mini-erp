package ports

import (
	"context"

	"github.com/isapr/mini-erp/services/api-gateway/internal/domain"
)

type BranchClient interface {
	CreateBranch(ctx context.Context, req CreateBranchRequest) (domain.BranchResponse, error)
	GetBranch(ctx context.Context, req GetBranchRequest) (domain.BranchResponse, error)
	ListBranches(ctx context.Context, req ListBranchesRequest) (domain.ListBranchesResponse, error)
	UpdateBranch(ctx context.Context, req UpdateBranchRequest) (domain.BranchResponse, error)
	CreateEmployeePlacement(ctx context.Context, req CreateEmployeePlacementRequest) (domain.PlacementResponse, error)
	ListAssignedBranches(ctx context.Context, req ListAssignedBranchesRequest) (ListAssignedBranchesResponse, error)
}

type CreateBranchRequest struct {
	BusinessID string `json:"business_id"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	Phone      string `json:"phone"`
}

type GetBranchRequest struct {
	BranchID string `json:"branch_id"`
}

type ListBranchesRequest struct {
	BusinessID string `json:"business_id"`
}

type UpdateBranchRequest struct {
	BranchID string `json:"branch_id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Status   string `json:"status"`
}

type CreateEmployeePlacementRequest struct {
	UserID         string `json:"user_id"`
	BusinessID     string `json:"business_id"`
	BranchID       string `json:"branch_id"`
	Position       string `json:"position"`
	EmploymentType string `json:"employment_type"`
}
