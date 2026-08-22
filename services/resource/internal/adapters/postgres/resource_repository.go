package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/isapr/mini-erp/services/resource/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ResourceRepository struct{ db *pgxpool.Pool }

func NewResourceRepository(db *pgxpool.Pool) *ResourceRepository { return &ResourceRepository{db: db} }

func (r *ResourceRepository) Create(ctx context.Context, resource domain.Resource) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO resources (id, business_id, branch_id, name, code, unit, type, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, resource.ID, resource.BusinessID, resource.BranchID, resource.Name, resource.Code, resource.Unit, resource.Type, resource.Status)
	return err
}

func (r *ResourceRepository) List(ctx context.Context, businessID string, branchID string) ([]domain.Resource, error) {
	query := `SELECT id, business_id, branch_id, name, code, unit, type, status FROM resources WHERE business_id = $1`
	args := []any{businessID}
	if branchID != "" {
		query += ` AND branch_id = $2`
		args = append(args, branchID)
	}
	query += ` ORDER BY name`
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	resources := []domain.Resource{}
	for rows.Next() {
		var resource domain.Resource
		if err := rows.Scan(&resource.ID, &resource.BusinessID, &resource.BranchID, &resource.Name, &resource.Code, &resource.Unit, &resource.Type, &resource.Status); err != nil {
			return nil, err
		}
		resources = append(resources, resource)
	}
	return resources, rows.Err()
}

func (r *ResourceRepository) FindByID(ctx context.Context, resourceID string) (domain.Resource, error) {
	var resource domain.Resource
	err := r.db.QueryRow(ctx, `
		SELECT id, business_id, branch_id, name, code, unit, type, status
		FROM resources
		WHERE id = $1
	`, resourceID).Scan(&resource.ID, &resource.BusinessID, &resource.BranchID, &resource.Name, &resource.Code, &resource.Unit, &resource.Type, &resource.Status)
	return resource, err
}

func (r *ResourceRepository) RecordStockMovement(ctx context.Context, movement domain.StockMovement) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO stock_movements (id, business_id, branch_id, resource_id, movement_type, quantity, reason, service_order_id, actor_user_id, request_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NULLIF($8, '00000000-0000-0000-0000-000000000000')::uuid, NULLIF($9, '00000000-0000-0000-0000-000000000000')::uuid, $10)
	`, movement.ID, movement.BusinessID, movement.BranchID, movement.ResourceID, movement.MovementType, movement.Quantity, movement.Reason, movement.ServiceOrderID, movement.ActorUserID, movement.RequestID)
	return err
}

func (r *ResourceRepository) Availability(ctx context.Context, resourceID string) (domain.ResourceAvailability, error) {
	var availability domain.ResourceAvailability
	availability.ResourceID = uuid.MustParse(resourceID)
	err := r.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(CASE WHEN movement_type = 'in' THEN quantity ELSE -quantity END), 0)
		FROM stock_movements
		WHERE resource_id = $1
	`, resourceID).Scan(&availability.Quantity)
	return availability, err
}

func (r *ResourceRepository) RecordResourceUsage(ctx context.Context, usage domain.ResourceUsage, movement domain.StockMovement) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		INSERT INTO stock_movements (id, business_id, branch_id, resource_id, movement_type, quantity, reason, service_order_id, actor_user_id, request_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NULLIF($9, '00000000-0000-0000-0000-000000000000')::uuid, $10)
	`, movement.ID, movement.BusinessID, movement.BranchID, movement.ResourceID, movement.MovementType, movement.Quantity, movement.Reason, movement.ServiceOrderID, movement.ActorUserID, movement.RequestID); err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO resource_usages (id, business_id, branch_id, service_order_id, resource_id, quantity, reason, recorded_by_user_id, stock_movement_id, request_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NULLIF($8, '00000000-0000-0000-0000-000000000000')::uuid, $9, $10)
	`, usage.ID, usage.BusinessID, usage.BranchID, usage.ServiceOrderID, usage.ResourceID, usage.Quantity, usage.Reason, usage.RecordedByUserID, usage.StockMovementID, usage.RequestID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *ResourceRepository) ListUsageByServiceOrder(ctx context.Context, serviceOrderID string) ([]domain.ResourceUsage, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, business_id, branch_id, service_order_id, resource_id, quantity, reason, recorded_by_user_id, stock_movement_id, request_id
		FROM resource_usages
		WHERE service_order_id = $1
		ORDER BY created_at DESC
	`, serviceOrderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usages := []domain.ResourceUsage{}
	for rows.Next() {
		var usage domain.ResourceUsage
		if err := rows.Scan(&usage.ID, &usage.BusinessID, &usage.BranchID, &usage.ServiceOrderID, &usage.ResourceID, &usage.Quantity, &usage.Reason, &usage.RecordedByUserID, &usage.StockMovementID, &usage.RequestID); err != nil {
			return nil, err
		}
		usages = append(usages, usage)
	}
	return usages, rows.Err()
}
