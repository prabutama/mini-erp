package application

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/isapr/mini-erp/services/resource/internal/domain"
	"github.com/isapr/mini-erp/services/resource/internal/ports"
)

var ErrValidation = errors.New("validation error")

type CreateResourceInput struct {
	BusinessID string
	BranchID   string
	Name       string
	Unit       string
	Type       string
}

type RecordStockMovementInput struct {
	BusinessID     string
	BranchID       string
	ResourceID     string
	MovementType   string
	Quantity       float64
	Reason         string
	ServiceOrderID string
	ActorUserID    string
	RequestID      string
}

type RecordResourceUsageInput struct {
	BusinessID       string
	BranchID         string
	ServiceOrderID   string
	ResourceID       string
	Quantity         float64
	Reason           string
	RecordedByUserID string
	RequestID        string
}

type ResourceService struct{ resources ports.ResourceRepository }

func NewResourceService(resources ports.ResourceRepository) *ResourceService {
	return &ResourceService{resources: resources}
}

func (s *ResourceService) Create(ctx context.Context, input CreateResourceInput) (domain.Resource, error) {
	businessID, err := uuid.Parse(input.BusinessID)
	if err != nil {
		return domain.Resource{}, ErrValidation
	}
	branchID, err := uuid.Parse(input.BranchID)
	if err != nil {
		return domain.Resource{}, ErrValidation
	}
	name := strings.TrimSpace(input.Name)
	unit := strings.TrimSpace(input.Unit)
	resourceType := strings.TrimSpace(input.Type)
	if name == "" || unit == "" || resourceType == "" {
		return domain.Resource{}, ErrValidation
	}
	resource := domain.Resource{ID: uuid.New(), BusinessID: businessID, BranchID: branchID, Name: name, Code: slug(name), Unit: unit, Type: resourceType, Status: "active"}
	if err := s.resources.Create(ctx, resource); err != nil {
		return domain.Resource{}, err
	}
	return resource, nil
}

func (s *ResourceService) List(ctx context.Context, businessID string, branchID string) ([]domain.Resource, error) {
	if _, err := uuid.Parse(businessID); err != nil {
		return nil, ErrValidation
	}
	if branchID != "" {
		if _, err := uuid.Parse(branchID); err != nil {
			return nil, ErrValidation
		}
	}
	return s.resources.List(ctx, businessID, branchID)
}

func (s *ResourceService) Get(ctx context.Context, resourceID string) (domain.Resource, error) {
	if _, err := uuid.Parse(resourceID); err != nil {
		return domain.Resource{}, ErrValidation
	}
	return s.resources.FindByID(ctx, resourceID)
}

func (s *ResourceService) RecordStockMovement(ctx context.Context, input RecordStockMovementInput) (domain.StockMovement, error) {
	businessID, err := uuid.Parse(input.BusinessID)
	if err != nil {
		return domain.StockMovement{}, ErrValidation
	}
	branchID, err := uuid.Parse(input.BranchID)
	if err != nil {
		return domain.StockMovement{}, ErrValidation
	}
	resourceID, err := uuid.Parse(input.ResourceID)
	if err != nil {
		return domain.StockMovement{}, ErrValidation
	}
	if input.Quantity <= 0 {
		return domain.StockMovement{}, ErrValidation
	}
	movementType := strings.TrimSpace(input.MovementType)
	if movementType != "in" && movementType != "out" {
		return domain.StockMovement{}, ErrValidation
	}
	movement := domain.StockMovement{ID: uuid.New(), BusinessID: businessID, BranchID: branchID, ResourceID: resourceID, MovementType: movementType, Quantity: input.Quantity, Reason: strings.TrimSpace(input.Reason), RequestID: input.RequestID}
	if input.ServiceOrderID != "" {
		id, err := uuid.Parse(input.ServiceOrderID)
		if err != nil {
			return domain.StockMovement{}, ErrValidation
		}
		movement.ServiceOrderID = id
	}
	if input.ActorUserID != "" {
		id, err := uuid.Parse(input.ActorUserID)
		if err != nil {
			return domain.StockMovement{}, ErrValidation
		}
		movement.ActorUserID = id
	}
	resource, err := s.resources.FindByID(ctx, input.ResourceID)
	if err != nil {
		return domain.StockMovement{}, err
	}
	if resource.BusinessID != businessID || resource.BranchID != branchID {
		return domain.StockMovement{}, ErrValidation
	}
	if err := s.resources.RecordStockMovement(ctx, movement); err != nil {
		return domain.StockMovement{}, err
	}
	return movement, nil
}

func (s *ResourceService) Availability(ctx context.Context, resourceID string) (domain.ResourceAvailability, error) {
	if _, err := uuid.Parse(resourceID); err != nil {
		return domain.ResourceAvailability{}, ErrValidation
	}
	return s.resources.Availability(ctx, resourceID)
}

func (s *ResourceService) RecordResourceUsage(ctx context.Context, input RecordResourceUsageInput) (domain.ResourceUsage, error) {
	businessID, err := uuid.Parse(input.BusinessID)
	if err != nil {
		return domain.ResourceUsage{}, ErrValidation
	}
	branchID, err := uuid.Parse(input.BranchID)
	if err != nil {
		return domain.ResourceUsage{}, ErrValidation
	}
	serviceOrderID, err := uuid.Parse(input.ServiceOrderID)
	if err != nil {
		return domain.ResourceUsage{}, ErrValidation
	}
	resourceID, err := uuid.Parse(input.ResourceID)
	if err != nil {
		return domain.ResourceUsage{}, ErrValidation
	}
	if input.Quantity <= 0 {
		return domain.ResourceUsage{}, ErrValidation
	}
	recordedByUserID := uuid.Nil
	if input.RecordedByUserID != "" {
		recordedByUserID, err = uuid.Parse(input.RecordedByUserID)
		if err != nil {
			return domain.ResourceUsage{}, ErrValidation
		}
	}
	resource, err := s.resources.FindByID(ctx, input.ResourceID)
	if err != nil {
		return domain.ResourceUsage{}, err
	}
	if resource.BusinessID != businessID || resource.BranchID != branchID {
		return domain.ResourceUsage{}, ErrValidation
	}
	movementID := uuid.New()
	usage := domain.ResourceUsage{ID: uuid.New(), BusinessID: businessID, BranchID: branchID, ServiceOrderID: serviceOrderID, ResourceID: resourceID, Quantity: input.Quantity, Reason: strings.TrimSpace(input.Reason), RecordedByUserID: recordedByUserID, StockMovementID: movementID, RequestID: input.RequestID}
	movement := domain.StockMovement{ID: movementID, BusinessID: businessID, BranchID: branchID, ResourceID: resourceID, MovementType: "out", Quantity: input.Quantity, Reason: usage.Reason, ServiceOrderID: serviceOrderID, ActorUserID: recordedByUserID, RequestID: input.RequestID}
	if err := s.resources.RecordResourceUsage(ctx, usage, movement); err != nil {
		return domain.ResourceUsage{}, err
	}
	return usage, nil
}

func (s *ResourceService) ListUsageByServiceOrder(ctx context.Context, serviceOrderID string) ([]domain.ResourceUsage, error) {
	if _, err := uuid.Parse(serviceOrderID); err != nil {
		return nil, ErrValidation
	}
	return s.resources.ListUsageByServiceOrder(ctx, serviceOrderID)
}

func slug(value string) string {
	clean := regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(strings.ToLower(strings.TrimSpace(value)), "-")
	return strings.Trim(clean, "-")
}
