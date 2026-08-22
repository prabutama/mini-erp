package ports

import (
	"context"

	"github.com/isapr/mini-erp/services/api-gateway/internal/domain"
)

type ReportingClient interface {
	GetAuditEvents(ctx context.Context, req GetAuditEventsRequest) (domain.ListAuditEventsResponse, error)
	GetOperationsSummary(ctx context.Context, req GetOperationsSummaryRequest) (domain.OperationsSummaryReportResponse, error)
}

type GetAuditEventsRequest struct {
	BusinessID string `json:"business_id"`
	BranchID   string `json:"branch_id"`
}
type GetOperationsSummaryRequest struct {
	BusinessID   string `json:"business_id"`
	BranchID     string `json:"branch_id"`
	SnapshotDate string `json:"snapshot_date"`
}
