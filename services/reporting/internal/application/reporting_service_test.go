package application

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/isapr/mini-erp/services/reporting/internal/domain"
)

type reportingRepositoryStub struct {
	event   domain.AuditEvent
	summary domain.OperationsSummary
}

func (s *reportingRepositoryStub) RecordAuditEvent(_ context.Context, event domain.AuditEvent) error {
	s.event = event
	return nil
}
func (s *reportingRepositoryStub) ListAuditEvents(context.Context, string, string) ([]domain.AuditEvent, error) {
	return []domain.AuditEvent{s.event}, nil
}
func (s *reportingRepositoryStub) UpsertOperationsSummary(_ context.Context, summary domain.OperationsSummary) error {
	s.summary = summary
	return nil
}
func (s *reportingRepositoryStub) GetOperationsSummary(context.Context, string, string, string) (domain.OperationsSummary, error) {
	return s.summary, nil
}

func TestRecordAuditEvent(t *testing.T) {
	repo := &reportingRepositoryStub{}
	service := NewReportingService(repo)
	event, err := service.RecordAuditEvent(context.Background(), RecordAuditEventInput{EventType: "service-order.created", EventVersion: 1, Producer: "operations-service", BusinessID: uuid.NewString(), Data: `{"ok":true}`})
	if err != nil {
		t.Fatal(err)
	}
	if event.ID == uuid.Nil || repo.event.EventType != "service-order.created" {
		t.Fatalf("unexpected event %#v", event)
	}
}

func TestUpsertOperationsSummary(t *testing.T) {
	repo := &reportingRepositoryStub{}
	service := NewReportingService(repo)
	summary, err := service.UpsertOperationsSummary(context.Background(), OperationsSummaryInput{BusinessID: uuid.NewString(), SnapshotDate: time.Now().UTC().Format("2006-01-02"), OpenOrders: 2})
	if err != nil {
		t.Fatal(err)
	}
	if summary.ID == uuid.Nil || summary.OpenOrders != 2 {
		t.Fatalf("unexpected summary %#v", summary)
	}
}
