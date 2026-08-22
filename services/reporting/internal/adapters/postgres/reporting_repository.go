package postgres

import (
	"context"

	"github.com/isapr/mini-erp/services/reporting/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReportingRepository struct{ db *pgxpool.Pool }

func NewReportingRepository(db *pgxpool.Pool) *ReportingRepository {
	return &ReportingRepository{db: db}
}

func (r *ReportingRepository) RecordAuditEvent(ctx context.Context, event domain.AuditEvent) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO audit_events (id, event_type, event_version, producer, business_id, branch_id, actor_id, request_id, occurred_at, data)
		VALUES ($1, $2, $3, $4, $5, NULLIF($6, '00000000-0000-0000-0000-000000000000')::uuid, NULLIF($7, '00000000-0000-0000-0000-000000000000')::uuid, $8, $9, $10::jsonb)
	`, event.ID, event.EventType, event.EventVersion, event.Producer, event.BusinessID, event.BranchID, event.ActorID, event.RequestID, event.OccurredAt, event.Data)
	return err
}

func (r *ReportingRepository) ListAuditEvents(ctx context.Context, businessID string, branchID string) ([]domain.AuditEvent, error) {
	query := `SELECT id, event_type, event_version, producer, business_id, COALESCE(branch_id, '00000000-0000-0000-0000-000000000000'::uuid), COALESCE(actor_id, '00000000-0000-0000-0000-000000000000'::uuid), COALESCE(request_id, ''), occurred_at, data::text FROM audit_events WHERE business_id = $1`
	args := []any{businessID}
	if branchID != "" {
		query += ` AND branch_id = $2`
		args = append(args, branchID)
	}
	query += ` ORDER BY occurred_at DESC LIMIT 100`
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	events := []domain.AuditEvent{}
	for rows.Next() {
		var event domain.AuditEvent
		if err := rows.Scan(&event.ID, &event.EventType, &event.EventVersion, &event.Producer, &event.BusinessID, &event.BranchID, &event.ActorID, &event.RequestID, &event.OccurredAt, &event.Data); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, rows.Err()
}

func (r *ReportingRepository) UpsertOperationsSummary(ctx context.Context, summary domain.OperationsSummary) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO operation_snapshots (id, business_id, branch_id, snapshot_date, open_orders, in_progress_orders, completed_orders, cancelled_orders, resources_used)
		VALUES ($1, $2, NULLIF($3, '00000000-0000-0000-0000-000000000000')::uuid, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (business_id, COALESCE(branch_id, '00000000-0000-0000-0000-000000000000'::uuid), snapshot_date) DO UPDATE SET
		open_orders = EXCLUDED.open_orders,
		in_progress_orders = EXCLUDED.in_progress_orders,
		completed_orders = EXCLUDED.completed_orders,
		cancelled_orders = EXCLUDED.cancelled_orders,
		resources_used = EXCLUDED.resources_used,
		updated_at = now()
	`, summary.ID, summary.BusinessID, summary.BranchID, summary.SnapshotDate, summary.OpenOrders, summary.InProgressOrders, summary.CompletedOrders, summary.CancelledOrders, summary.ResourcesUsed)
	return err
}

func (r *ReportingRepository) GetOperationsSummary(ctx context.Context, businessID string, branchID string, snapshotDate string) (domain.OperationsSummary, error) {
	query := `SELECT id, business_id, COALESCE(branch_id, '00000000-0000-0000-0000-000000000000'::uuid), snapshot_date, open_orders, in_progress_orders, completed_orders, cancelled_orders, resources_used FROM operation_snapshots WHERE business_id = $1 AND snapshot_date = $2`
	args := []any{businessID, snapshotDate}
	if branchID != "" {
		query += ` AND branch_id = $3`
		args = append(args, branchID)
	}
	query += ` ORDER BY updated_at DESC LIMIT 1`
	var summary domain.OperationsSummary
	err := r.db.QueryRow(ctx, query, args...).Scan(&summary.ID, &summary.BusinessID, &summary.BranchID, &summary.SnapshotDate, &summary.OpenOrders, &summary.InProgressOrders, &summary.CompletedOrders, &summary.CancelledOrders, &summary.ResourcesUsed)
	return summary, err
}
