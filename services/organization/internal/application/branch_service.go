package application

import (
	"context"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/isapr/mini-erp/services/organization/internal/domain"
	"github.com/isapr/mini-erp/services/organization/internal/ports"
)

type CreateBranchInput struct {
	BusinessID string
	Name       string
	Address    string
	Phone      string
}

type UpdateBranchInput struct {
	BranchID string
	Name     string
	Address  string
	Phone    string
	Status   string
}

type BranchService struct {
	branches ports.BranchRepository
}

func NewBranchService(branches ports.BranchRepository) *BranchService {
	return &BranchService{branches: branches}
}

func (s *BranchService) CreateBranch(ctx context.Context, input CreateBranchInput) (domain.Branch, error) {
	businessID, err := uuid.Parse(input.BusinessID)
	if err != nil || input.Name == "" {
		return domain.Branch{}, ErrValidation
	}

	branch := domain.Branch{
		ID:         uuid.New(),
		BusinessID: businessID,
		Name:       strings.TrimSpace(input.Name),
		Code:       branchCode(input.Name),
		Address:    strings.TrimSpace(input.Address),
		Phone:      strings.TrimSpace(input.Phone),
		Status:     "active",
	}

	if err := s.branches.Create(ctx, branch); err != nil {
		return domain.Branch{}, err
	}
	return branch, nil
}

func (s *BranchService) GetBranch(ctx context.Context, branchID string) (domain.Branch, error) {
	return s.branches.FindByID(ctx, branchID)
}

func (s *BranchService) ListBranches(ctx context.Context, businessID string) ([]domain.Branch, error) {
	if _, err := uuid.Parse(businessID); err != nil {
		return nil, ErrValidation
	}
	return s.branches.ListByBusiness(ctx, businessID)
}

func (s *BranchService) UpdateBranch(ctx context.Context, input UpdateBranchInput) (domain.Branch, error) {
	branch, err := s.branches.FindByID(ctx, input.BranchID)
	if err != nil {
		return domain.Branch{}, err
	}
	if input.Name != "" {
		branch.Name = strings.TrimSpace(input.Name)
	}
	if input.Address != "" {
		branch.Address = strings.TrimSpace(input.Address)
	}
	if input.Phone != "" {
		branch.Phone = strings.TrimSpace(input.Phone)
	}
	if input.Status != "" {
		branch.Status = strings.TrimSpace(input.Status)
	}
	if err := s.branches.Update(ctx, branch); err != nil {
		return domain.Branch{}, err
	}
	return branch, nil
}

func branchCode(name string) string {
	code := strings.ToLower(strings.TrimSpace(name))
	code = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(code, "-")
	code = strings.Trim(code, "-")
	if code == "" {
		return uuid.NewString()
	}
	return code
}
