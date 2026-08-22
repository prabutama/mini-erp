package ports

import (
	"context"

	"github.com/isapr/mini-erp/services/reporting/internal/domain"
)

type ReportingRepository interface {
	RecordAuditEvent(ctx context.Context, event domain.AuditEvent) error
	ListAuditEvents(ctx context.Context, businessID string, branchID string) ([]domain.AuditEvent, error)
	UpsertOperationsSummary(ctx context.Context, summary domain.OperationsSummary) error
	GetOperationsSummary(ctx context.Context, businessID string, branchID string, snapshotDate string) (domain.OperationsSummary, error)
}
