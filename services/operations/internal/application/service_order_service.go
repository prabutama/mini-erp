package application

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/isapr/mini-erp/services/operations/internal/domain"
	"github.com/isapr/mini-erp/services/operations/internal/ports"
)

var ErrInvalidStatusTransition = errors.New("invalid status transition")
var ErrServiceOrderClosed = errors.New("service order closed")

type CreateServiceOrderInput struct {
	BusinessID          string
	BranchID            string
	ServiceDefinitionID string
	Title               string
	Description         string
	Priority            string
}

type TransitionServiceOrderInput struct {
	ServiceOrderID  string
	BusinessID      string
	Status          string
	ChangedByUserID string
	RequestID       string
}

type AssignServiceOrderInput struct {
	ServiceOrderID   string
	BusinessID       string
	AssignedUserID   string
	AssignedByUserID string
	RequestID        string
}

type ListServiceOrdersInput struct {
	BusinessID     string
	BranchID       string
	Status         string
	AssignedUserID string
}

type ServiceOrderSummaryInput struct {
	BusinessID     string
	BranchID       string
	AssignedUserID string
}

type ServiceOrderService struct {
	orders    ports.ServiceOrderRepository
	services  ports.ServiceDefinitionRepository
	publisher ports.EventPublisher
}

func NewServiceOrderService(orders ports.ServiceOrderRepository, services ports.ServiceDefinitionRepository, publisher ...ports.EventPublisher) *ServiceOrderService {
	var eventPublisher ports.EventPublisher
	if len(publisher) > 0 {
		eventPublisher = publisher[0]
	}
	return &ServiceOrderService{orders: orders, services: services, publisher: eventPublisher}
}

func (s *ServiceOrderService) Create(ctx context.Context, input CreateServiceOrderInput) (domain.ServiceOrder, error) {
	businessID, err := uuid.Parse(input.BusinessID)
	if err != nil {
		return domain.ServiceOrder{}, ErrValidation
	}
	branchID, err := uuid.Parse(input.BranchID)
	if err != nil {
		return domain.ServiceOrder{}, ErrValidation
	}
	serviceDefinitionID, err := uuid.Parse(input.ServiceDefinitionID)
	if err != nil {
		return domain.ServiceOrder{}, ErrValidation
	}
	serviceDefinition, err := s.services.FindByID(ctx, input.ServiceDefinitionID)
	if err != nil {
		return domain.ServiceOrder{}, err
	}
	if serviceDefinition.BusinessID != businessID {
		return domain.ServiceOrder{}, ErrValidation
	}
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return domain.ServiceOrder{}, ErrValidation
	}
	priority := strings.TrimSpace(input.Priority)
	if priority == "" {
		priority = "normal"
	}
	order := domain.ServiceOrder{ID: uuid.New(), BusinessID: businessID, BranchID: branchID, ServiceDefinitionID: serviceDefinitionID, Title: title, Description: strings.TrimSpace(input.Description), Status: "open", Priority: priority}
	if err := s.orders.Create(ctx, order); err != nil {
		return domain.ServiceOrder{}, err
	}
	s.publish(ctx, domain.EventEnvelope{EventID: uuid.NewString(), EventType: "service-order.created", EventVersion: 1, OccurredAt: time.Now().UTC(), Producer: "operations-service", BusinessID: order.BusinessID.String(), BranchID: order.BranchID.String(), Data: map[string]any{"service_order_id": order.ID.String(), "service_definition_id": order.ServiceDefinitionID.String(), "title": order.Title, "status": order.Status, "priority": order.Priority}})
	return order, nil
}

func (s *ServiceOrderService) List(ctx context.Context, input ListServiceOrdersInput) ([]domain.ServiceOrder, error) {
	if _, err := uuid.Parse(input.BusinessID); err != nil {
		return nil, ErrValidation
	}
	if input.BranchID != "" {
		if _, err := uuid.Parse(input.BranchID); err != nil {
			return nil, ErrValidation
		}
	}
	if input.AssignedUserID != "" {
		if _, err := uuid.Parse(input.AssignedUserID); err != nil {
			return nil, ErrValidation
		}
	}
	return s.orders.ListByBusiness(ctx, input.BusinessID, input.BranchID, strings.TrimSpace(input.Status), input.AssignedUserID)
}

func (s *ServiceOrderService) Summary(ctx context.Context, input ServiceOrderSummaryInput) (domain.ServiceOrderSummary, error) {
	if _, err := uuid.Parse(input.BusinessID); err != nil {
		return domain.ServiceOrderSummary{}, ErrValidation
	}
	if input.BranchID != "" {
		if _, err := uuid.Parse(input.BranchID); err != nil {
			return domain.ServiceOrderSummary{}, ErrValidation
		}
	}
	if input.AssignedUserID != "" {
		if _, err := uuid.Parse(input.AssignedUserID); err != nil {
			return domain.ServiceOrderSummary{}, ErrValidation
		}
	}
	return s.orders.SummaryByBusiness(ctx, input.BusinessID, input.BranchID, input.AssignedUserID)
}

func (s *ServiceOrderService) Get(ctx context.Context, serviceOrderID string) (domain.ServiceOrder, error) {
	if _, err := uuid.Parse(serviceOrderID); err != nil {
		return domain.ServiceOrder{}, ErrValidation
	}
	return s.orders.FindByID(ctx, serviceOrderID)
}

func (s *ServiceOrderService) Transition(ctx context.Context, input TransitionServiceOrderInput) (domain.ServiceOrder, error) {
	if _, err := uuid.Parse(input.ServiceOrderID); err != nil {
		return domain.ServiceOrder{}, ErrValidation
	}
	businessID, err := uuid.Parse(input.BusinessID)
	if err != nil {
		return domain.ServiceOrder{}, ErrValidation
	}
	if input.ChangedByUserID != "" {
		if _, err := uuid.Parse(input.ChangedByUserID); err != nil {
			return domain.ServiceOrder{}, ErrValidation
		}
	}

	order, err := s.orders.FindByID(ctx, input.ServiceOrderID)
	if err != nil {
		return domain.ServiceOrder{}, err
	}
	if order.BusinessID != businessID {
		return domain.ServiceOrder{}, ErrValidation
	}
	nextStatus := strings.TrimSpace(input.Status)
	if !canTransition(order.Status, nextStatus) {
		return domain.ServiceOrder{}, ErrInvalidStatusTransition
	}
	previousStatus := order.Status
	order.Status = nextStatus
	if err := s.orders.UpdateStatus(ctx, order, previousStatus, input.ChangedByUserID, input.RequestID); err != nil {
		return domain.ServiceOrder{}, err
	}
	s.publish(ctx, domain.EventEnvelope{EventID: uuid.NewString(), EventType: "service-order.status-changed", EventVersion: 1, OccurredAt: time.Now().UTC(), Producer: "operations-service", BusinessID: order.BusinessID.String(), BranchID: order.BranchID.String(), ActorID: input.ChangedByUserID, RequestID: input.RequestID, Data: map[string]any{"service_order_id": order.ID.String(), "from_status": previousStatus, "to_status": order.Status}})
	return order, nil
}

func (s *ServiceOrderService) Assign(ctx context.Context, input AssignServiceOrderInput) (domain.ServiceOrderAssignment, error) {
	serviceOrderID, err := uuid.Parse(input.ServiceOrderID)
	if err != nil {
		return domain.ServiceOrderAssignment{}, ErrValidation
	}
	businessID, err := uuid.Parse(input.BusinessID)
	if err != nil {
		return domain.ServiceOrderAssignment{}, ErrValidation
	}
	assignedUserID, err := uuid.Parse(input.AssignedUserID)
	if err != nil {
		return domain.ServiceOrderAssignment{}, ErrValidation
	}
	assignedByUserID, err := uuid.Parse(input.AssignedByUserID)
	if err != nil {
		return domain.ServiceOrderAssignment{}, ErrValidation
	}

	order, err := s.orders.FindByID(ctx, input.ServiceOrderID)
	if err != nil {
		return domain.ServiceOrderAssignment{}, err
	}
	if order.BusinessID != businessID {
		return domain.ServiceOrderAssignment{}, ErrValidation
	}
	if order.Status == "completed" || order.Status == "cancelled" {
		return domain.ServiceOrderAssignment{}, ErrServiceOrderClosed
	}

	assignment := domain.ServiceOrderAssignment{ID: uuid.New(), ServiceOrderID: serviceOrderID, BusinessID: businessID, BranchID: order.BranchID, AssignedUserID: assignedUserID, AssignedByUserID: assignedByUserID, Status: "active", RequestID: input.RequestID}
	if err := s.orders.Assign(ctx, assignment); err != nil {
		return domain.ServiceOrderAssignment{}, err
	}
	s.publish(ctx, domain.EventEnvelope{EventID: uuid.NewString(), EventType: "service-order.assigned", EventVersion: 1, OccurredAt: time.Now().UTC(), Producer: "operations-service", BusinessID: assignment.BusinessID.String(), BranchID: assignment.BranchID.String(), ActorID: assignment.AssignedByUserID.String(), RequestID: assignment.RequestID, Data: map[string]any{"assignment_id": assignment.ID.String(), "service_order_id": assignment.ServiceOrderID.String(), "assigned_user_id": assignment.AssignedUserID.String(), "status": assignment.Status}})
	return assignment, nil
}

func (s *ServiceOrderService) publish(ctx context.Context, event domain.EventEnvelope) {
	if s.publisher == nil {
		return
	}
	if err := s.publisher.Publish(ctx, event); err != nil {
		log.Printf("publish %s failed: %v", event.EventType, err)
	}
}

func (s *ServiceOrderService) ListAssignedToUser(ctx context.Context, businessID string, assignedUserID string, branchID string) ([]domain.ServiceOrder, error) {
	if _, err := uuid.Parse(businessID); err != nil {
		return nil, ErrValidation
	}
	if _, err := uuid.Parse(assignedUserID); err != nil {
		return nil, ErrValidation
	}
	if branchID != "" {
		if _, err := uuid.Parse(branchID); err != nil {
			return nil, ErrValidation
		}
	}
	return s.orders.ListAssignedToUser(ctx, businessID, assignedUserID, branchID)
}

func (s *ServiceOrderService) ListAssignments(ctx context.Context, serviceOrderID string) ([]domain.ServiceOrderAssignment, error) {
	if _, err := uuid.Parse(serviceOrderID); err != nil {
		return nil, ErrValidation
	}
	return s.orders.ListAssignmentsByOrder(ctx, serviceOrderID)
}

func canTransition(from string, to string) bool {
	allowed := map[string]map[string]struct{}{
		"open":        {"in_progress": {}, "cancelled": {}},
		"in_progress": {"completed": {}, "cancelled": {}},
	}
	statuses, ok := allowed[from]
	if !ok {
		return false
	}
	_, ok = statuses[to]
	return ok
}
