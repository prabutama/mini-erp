package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/isapr/mini-erp/services/reporting/internal/domain"
	"github.com/isapr/mini-erp/services/reporting/internal/ports"
)

var ErrValidation = errors.New("validation error")

type RecordAuditEventInput struct {
	EventType    string
	EventVersion int
	Producer     string
	BusinessID   string
	BranchID     string
	ActorID      string
	RequestID    string
	OccurredAt   time.Time
	Data         string
}
type OperationsSummaryInput struct {
	BusinessID       string
	BranchID         string
	SnapshotDate     string
	OpenOrders       int
	InProgressOrders int
	CompletedOrders  int
	CancelledOrders  int
	ResourcesUsed    float64
}

type ReportingService struct{ reports ports.ReportingRepository }

func NewReportingService(reports ports.ReportingRepository) *ReportingService {
	return &ReportingService{reports: reports}
}

func (s *ReportingService) RecordAuditEvent(ctx context.Context, input RecordAuditEventInput) (domain.AuditEvent, error) {
	businessID, err := uuid.Parse(input.BusinessID)
	if err != nil {
		return domain.AuditEvent{}, ErrValidation
	}
	branchID := uuid.Nil
	if input.BranchID != "" {
		branchID, err = uuid.Parse(input.BranchID)
		if err != nil {
			return domain.AuditEvent{}, ErrValidation
		}
	}
	actorID := uuid.Nil
	if input.ActorID != "" {
		actorID, err = uuid.Parse(input.ActorID)
		if err != nil {
			return domain.AuditEvent{}, ErrValidation
		}
	}
	if strings.TrimSpace(input.EventType) == "" || strings.TrimSpace(input.Producer) == "" || input.EventVersion <= 0 {
		return domain.AuditEvent{}, ErrValidation
	}
	data := strings.TrimSpace(input.Data)
	if data == "" {
		data = "{}"
	}
	occurredAt := input.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = time.Now().UTC()
	}
	event := domain.AuditEvent{ID: uuid.New(), EventType: input.EventType, EventVersion: input.EventVersion, Producer: input.Producer, BusinessID: businessID, BranchID: branchID, ActorID: actorID, RequestID: input.RequestID, OccurredAt: occurredAt, Data: data}
	if err := s.reports.RecordAuditEvent(ctx, event); err != nil {
		return domain.AuditEvent{}, err
	}
	return event, nil
}

func (s *ReportingService) ListAuditEvents(ctx context.Context, businessID string, branchID string) ([]domain.AuditEvent, error) {
	if _, err := uuid.Parse(businessID); err != nil {
		return nil, ErrValidation
	}
	if branchID != "" {
		if _, err := uuid.Parse(branchID); err != nil {
			return nil, ErrValidation
		}
	}
	return s.reports.ListAuditEvents(ctx, businessID, branchID)
}

func (s *ReportingService) UpsertOperationsSummary(ctx context.Context, input OperationsSummaryInput) (domain.OperationsSummary, error) {
	businessID, err := uuid.Parse(input.BusinessID)
	if err != nil {
		return domain.OperationsSummary{}, ErrValidation
	}
	branchID := uuid.Nil
	if input.BranchID != "" {
		branchID, err = uuid.Parse(input.BranchID)
		if err != nil {
			return domain.OperationsSummary{}, ErrValidation
		}
	}
	snapshotDate, err := time.Parse("2006-01-02", input.SnapshotDate)
	if err != nil {
		return domain.OperationsSummary{}, ErrValidation
	}
	summary := domain.OperationsSummary{ID: uuid.New(), BusinessID: businessID, BranchID: branchID, SnapshotDate: snapshotDate, OpenOrders: input.OpenOrders, InProgressOrders: input.InProgressOrders, CompletedOrders: input.CompletedOrders, CancelledOrders: input.CancelledOrders, ResourcesUsed: input.ResourcesUsed}
	if err := s.reports.UpsertOperationsSummary(ctx, summary); err != nil {
		return domain.OperationsSummary{}, err
	}
	return summary, nil
}

func (s *ReportingService) GetOperationsSummary(ctx context.Context, businessID string, branchID string, snapshotDate string) (domain.OperationsSummary, error) {
	if _, err := uuid.Parse(businessID); err != nil {
		return domain.OperationsSummary{}, ErrValidation
	}
	if branchID != "" {
		if _, err := uuid.Parse(branchID); err != nil {
			return domain.OperationsSummary{}, ErrValidation
		}
	}
	if snapshotDate == "" {
		snapshotDate = time.Now().UTC().Format("2006-01-02")
	}
	if _, err := time.Parse("2006-01-02", snapshotDate); err != nil {
		return domain.OperationsSummary{}, ErrValidation
	}
	return s.reports.GetOperationsSummary(ctx, businessID, branchID, snapshotDate)
}
