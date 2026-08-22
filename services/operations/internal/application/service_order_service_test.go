package application

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/isapr/mini-erp/services/operations/internal/domain"
)

type serviceDefinitionRepositoryStub struct{ service domain.ServiceDefinition }

func (s serviceDefinitionRepositoryStub) Create(context.Context, domain.ServiceDefinition) error {
	return nil
}

func (s serviceDefinitionRepositoryStub) ListByBusiness(context.Context, string) ([]domain.ServiceDefinition, error) {
	return []domain.ServiceDefinition{s.service}, nil
}

func (s serviceDefinitionRepositoryStub) FindByID(context.Context, string) (domain.ServiceDefinition, error) {
	return s.service, nil
}

type serviceOrderRepositoryStub struct{ order domain.ServiceOrder }

func (s *serviceOrderRepositoryStub) Create(_ context.Context, order domain.ServiceOrder) error {
	s.order = order
	return nil
}

func (s *serviceOrderRepositoryStub) ListByBusiness(context.Context, string, string, string, string) ([]domain.ServiceOrder, error) {
	return []domain.ServiceOrder{s.order}, nil
}

func (s *serviceOrderRepositoryStub) SummaryByBusiness(context.Context, string, string, string) (domain.ServiceOrderSummary, error) {
	return domain.ServiceOrderSummary{Total: 1, Open: 1}, nil
}

func (s *serviceOrderRepositoryStub) FindByID(context.Context, string) (domain.ServiceOrder, error) {
	return s.order, nil
}

func (s *serviceOrderRepositoryStub) UpdateStatus(_ context.Context, order domain.ServiceOrder, _ string, _ string, _ string) error {
	s.order = order
	return nil
}

func (s *serviceOrderRepositoryStub) Assign(context.Context, domain.ServiceOrderAssignment) error {
	return nil
}

func (s *serviceOrderRepositoryStub) ListAssignedToUser(context.Context, string, string, string) ([]domain.ServiceOrder, error) {
	return []domain.ServiceOrder{s.order}, nil
}

func (s *serviceOrderRepositoryStub) ListAssignmentsByOrder(context.Context, string) ([]domain.ServiceOrderAssignment, error) {
	return []domain.ServiceOrderAssignment{{ID: uuid.New(), ServiceOrderID: s.order.ID, BusinessID: s.order.BusinessID, BranchID: s.order.BranchID, Status: "active"}}, nil
}

func TestCreateServiceOrder(t *testing.T) {
	businessID := uuid.New()
	branchID := uuid.New()
	serviceDefinitionID := uuid.New()
	orders := &serviceOrderRepositoryStub{}
	service := NewServiceOrderService(orders, serviceDefinitionRepositoryStub{service: domain.ServiceDefinition{ID: serviceDefinitionID, BusinessID: businessID, Status: "active"}})

	order, err := service.Create(context.Background(), CreateServiceOrderInput{BusinessID: businessID.String(), BranchID: branchID.String(), ServiceDefinitionID: serviceDefinitionID.String(), Title: "Fix AC"})
	if err != nil {
		t.Fatal(err)
	}
	if order.Status != "open" || order.Priority != "normal" || orders.order.ID == uuid.Nil {
		t.Fatalf("unexpected order %#v", order)
	}
}

func TestCreateServiceOrderRejectsOtherBusinessServiceDefinition(t *testing.T) {
	service := NewServiceOrderService(&serviceOrderRepositoryStub{}, serviceDefinitionRepositoryStub{service: domain.ServiceDefinition{ID: uuid.New(), BusinessID: uuid.New(), Status: "active"}})

	_, err := service.Create(context.Background(), CreateServiceOrderInput{BusinessID: uuid.NewString(), BranchID: uuid.NewString(), ServiceDefinitionID: uuid.NewString(), Title: "Fix AC"})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestTransitionServiceOrder(t *testing.T) {
	businessID := uuid.New()
	orderID := uuid.New()
	orders := &serviceOrderRepositoryStub{order: domain.ServiceOrder{ID: orderID, BusinessID: businessID, BranchID: uuid.New(), Status: "open"}}
	service := NewServiceOrderService(orders, serviceDefinitionRepositoryStub{})

	order, err := service.Transition(context.Background(), TransitionServiceOrderInput{ServiceOrderID: orderID.String(), BusinessID: businessID.String(), Status: "in_progress", ChangedByUserID: uuid.NewString(), RequestID: "request-1"})
	if err != nil {
		t.Fatal(err)
	}
	if order.Status != "in_progress" || orders.order.Status != "in_progress" {
		t.Fatalf("expected in_progress, got %#v", order)
	}
}

func TestTransitionServiceOrderRejectsInvalidTransition(t *testing.T) {
	businessID := uuid.New()
	orderID := uuid.New()
	orders := &serviceOrderRepositoryStub{order: domain.ServiceOrder{ID: orderID, BusinessID: businessID, BranchID: uuid.New(), Status: "open"}}
	service := NewServiceOrderService(orders, serviceDefinitionRepositoryStub{})

	_, err := service.Transition(context.Background(), TransitionServiceOrderInput{ServiceOrderID: orderID.String(), BusinessID: businessID.String(), Status: "completed"})
	if !errors.Is(err, ErrInvalidStatusTransition) {
		t.Fatalf("expected invalid transition error, got %v", err)
	}
}

func TestAssignServiceOrder(t *testing.T) {
	businessID := uuid.New()
	orderID := uuid.New()
	orders := &serviceOrderRepositoryStub{order: domain.ServiceOrder{ID: orderID, BusinessID: businessID, BranchID: uuid.New(), Status: "open"}}
	service := NewServiceOrderService(orders, serviceDefinitionRepositoryStub{})

	assignment, err := service.Assign(context.Background(), AssignServiceOrderInput{ServiceOrderID: orderID.String(), BusinessID: businessID.String(), AssignedUserID: uuid.NewString(), AssignedByUserID: uuid.NewString(), RequestID: "request-1"})
	if err != nil {
		t.Fatal(err)
	}
	if assignment.ID == uuid.Nil || assignment.Status != "active" {
		t.Fatalf("unexpected assignment %#v", assignment)
	}
}

func TestAssignServiceOrderRejectsCompletedOrder(t *testing.T) {
	businessID := uuid.New()
	orderID := uuid.New()
	orders := &serviceOrderRepositoryStub{order: domain.ServiceOrder{ID: orderID, BusinessID: businessID, BranchID: uuid.New(), Status: "completed"}}
	service := NewServiceOrderService(orders, serviceDefinitionRepositoryStub{})

	_, err := service.Assign(context.Background(), AssignServiceOrderInput{ServiceOrderID: orderID.String(), BusinessID: businessID.String(), AssignedUserID: uuid.NewString(), AssignedByUserID: uuid.NewString()})
	if !errors.Is(err, ErrServiceOrderClosed) {
		t.Fatalf("expected service order closed, got %v", err)
	}
}

func TestListAssignedToUser(t *testing.T) {
	businessID := uuid.New()
	userID := uuid.New()
	orders := &serviceOrderRepositoryStub{order: domain.ServiceOrder{ID: uuid.New(), BusinessID: businessID, BranchID: uuid.New(), Status: "open"}}
	service := NewServiceOrderService(orders, serviceDefinitionRepositoryStub{})

	result, err := service.ListAssignedToUser(context.Background(), businessID.String(), userID.String(), "")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("expected one order, got %d", len(result))
	}
}
