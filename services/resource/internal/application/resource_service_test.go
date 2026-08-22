package application

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/isapr/mini-erp/services/resource/internal/domain"
)

type resourceRepositoryStub struct {
	resource domain.Resource
	movement domain.StockMovement
}

func (s *resourceRepositoryStub) Create(_ context.Context, resource domain.Resource) error {
	s.resource = resource
	return nil
}
func (s *resourceRepositoryStub) List(context.Context, string, string) ([]domain.Resource, error) {
	return []domain.Resource{s.resource}, nil
}
func (s *resourceRepositoryStub) FindByID(context.Context, string) (domain.Resource, error) {
	return s.resource, nil
}
func (s *resourceRepositoryStub) RecordStockMovement(_ context.Context, movement domain.StockMovement) error {
	s.movement = movement
	return nil
}
func (s *resourceRepositoryStub) Availability(context.Context, string) (domain.ResourceAvailability, error) {
	return domain.ResourceAvailability{ResourceID: s.resource.ID, Quantity: 10}, nil
}
func (s *resourceRepositoryStub) RecordResourceUsage(_ context.Context, usage domain.ResourceUsage, movement domain.StockMovement) error {
	s.movement = movement
	return nil
}
func (s *resourceRepositoryStub) ListUsageByServiceOrder(context.Context, string) ([]domain.ResourceUsage, error) {
	return []domain.ResourceUsage{{ID: uuid.New(), BusinessID: s.resource.BusinessID, BranchID: s.resource.BranchID, ResourceID: s.resource.ID, Quantity: 3}}, nil
}

func TestCreateResource(t *testing.T) {
	repo := &resourceRepositoryStub{}
	service := NewResourceService(repo)
	resource, err := service.Create(context.Background(), CreateResourceInput{BusinessID: uuid.NewString(), BranchID: uuid.NewString(), Name: "Filter", Unit: "pcs", Type: "stock"})
	if err != nil {
		t.Fatal(err)
	}
	if resource.ID == uuid.Nil || resource.Code != "filter" || repo.resource.ID == uuid.Nil {
		t.Fatalf("unexpected resource %#v", resource)
	}
}

func TestRecordStockMovement(t *testing.T) {
	businessID := uuid.New()
	branchID := uuid.New()
	resourceID := uuid.New()
	repo := &resourceRepositoryStub{resource: domain.Resource{ID: resourceID, BusinessID: businessID, BranchID: branchID}}
	service := NewResourceService(repo)
	movement, err := service.RecordStockMovement(context.Background(), RecordStockMovementInput{BusinessID: businessID.String(), BranchID: branchID.String(), ResourceID: resourceID.String(), MovementType: "in", Quantity: 5})
	if err != nil {
		t.Fatal(err)
	}
	if movement.ID == uuid.Nil || movement.Quantity != 5 {
		t.Fatalf("unexpected movement %#v", movement)
	}
}

func TestRecordResourceUsage(t *testing.T) {
	businessID := uuid.New()
	branchID := uuid.New()
	resourceID := uuid.New()
	serviceOrderID := uuid.New()
	repo := &resourceRepositoryStub{resource: domain.Resource{ID: resourceID, BusinessID: businessID, BranchID: branchID}}
	service := NewResourceService(repo)
	usage, err := service.RecordResourceUsage(context.Background(), RecordResourceUsageInput{BusinessID: businessID.String(), BranchID: branchID.String(), ServiceOrderID: serviceOrderID.String(), ResourceID: resourceID.String(), Quantity: 3})
	if err != nil {
		t.Fatal(err)
	}
	if usage.ID == uuid.Nil || usage.StockMovementID == uuid.Nil || repo.movement.MovementType != "out" {
		t.Fatalf("unexpected usage %#v movement %#v", usage, repo.movement)
	}
}
