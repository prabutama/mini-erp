package application

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/isapr/mini-erp/services/organization/internal/domain"
	"github.com/isapr/mini-erp/services/organization/internal/ports"
)

type CreatePlacementInput struct {
	UserID         string
	BusinessID     string
	BranchID       string
	Position       string
	EmploymentType string
}

type ListAssignedBranchesInput struct {
	UserID     string
	BusinessID string
}

type PlacementService struct {
	placements ports.PlacementRepository
}

func NewPlacementService(placements ports.PlacementRepository) *PlacementService {
	return &PlacementService{placements: placements}
}

func (s *PlacementService) CreatePlacement(ctx context.Context, input CreatePlacementInput) (domain.EmployeePlacement, error) {
	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		return domain.EmployeePlacement{}, ErrValidation
	}
	businessID, err := uuid.Parse(input.BusinessID)
	if err != nil {
		return domain.EmployeePlacement{}, ErrValidation
	}
	branchID, err := uuid.Parse(input.BranchID)
	if err != nil {
		return domain.EmployeePlacement{}, ErrValidation
	}
	placement := domain.EmployeePlacement{ID: uuid.New(), UserID: userID, BusinessID: businessID, BranchID: branchID, Position: strings.TrimSpace(input.Position), EmploymentType: strings.TrimSpace(input.EmploymentType), Status: "active", StartDate: time.Now()}
	if placement.Position == "" || placement.EmploymentType == "" {
		return domain.EmployeePlacement{}, ErrValidation
	}
	if err := s.placements.Create(ctx, placement); err != nil {
		return domain.EmployeePlacement{}, err
	}
	return placement, nil
}

func (s *PlacementService) ListAssignedBranches(ctx context.Context, input ListAssignedBranchesInput) ([]domain.EmployeePlacement, error) {
	if _, err := uuid.Parse(input.UserID); err != nil {
		return nil, ErrValidation
	}
	if _, err := uuid.Parse(input.BusinessID); err != nil {
		return nil, ErrValidation
	}
	return s.placements.ListAssignedBranches(ctx, input.UserID, input.BusinessID)
}
