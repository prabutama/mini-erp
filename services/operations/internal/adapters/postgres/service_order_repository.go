package postgres

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"github.com/isapr/mini-erp/services/operations/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceOrderRepository struct {
	db *pgxpool.Pool
}

func NewServiceOrderRepository(db *pgxpool.Pool) *ServiceOrderRepository {
	return &ServiceOrderRepository{db: db}
}

func (r *ServiceOrderRepository) Create(ctx context.Context, order domain.ServiceOrder) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO service_orders (id, business_id, branch_id, service_definition_id, title, description, status, priority)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, order.ID, order.BusinessID, order.BranchID, order.ServiceDefinitionID, order.Title, order.Description, order.Status, order.Priority)
	return err
}

func (r *ServiceOrderRepository) ListByBusiness(ctx context.Context, businessID string, branchID string, status string, assignedUserID string) ([]domain.ServiceOrder, error) {
	query := `
		SELECT DISTINCT o.id, o.business_id, o.branch_id, o.service_definition_id, o.title, o.description, o.status, o.priority, o.created_at
		FROM service_orders o`
	args := []any{businessID}
	if assignedUserID != "" {
		query += ` JOIN service_order_assignments a ON a.service_order_id = o.id AND a.status = 'active'`
	}
	query += ` WHERE o.business_id = $1`
	if branchID != "" {
		args = append(args, branchID)
		query += ` AND o.branch_id = $` + strconv.Itoa(len(args))
	}
	if status != "" {
		args = append(args, branchID)
		args[len(args)-1] = status
		query += ` AND o.status = $` + strconv.Itoa(len(args))
	}
	if assignedUserID != "" {
		args = append(args, assignedUserID)
		query += ` AND a.assigned_user_id = $` + strconv.Itoa(len(args))
	}
	query += ` ORDER BY o.created_at DESC`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []domain.ServiceOrder{}
	for rows.Next() {
		var order domain.ServiceOrder
		var ignored any
		if err := rows.Scan(&order.ID, &order.BusinessID, &order.BranchID, &order.ServiceDefinitionID, &order.Title, &order.Description, &order.Status, &order.Priority, &ignored); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, rows.Err()
}

func (r *ServiceOrderRepository) SummaryByBusiness(ctx context.Context, businessID string, branchID string, assignedUserID string) (domain.ServiceOrderSummary, error) {
	query := `
		SELECT o.status, count(*)
		FROM service_orders o`
	args := []any{businessID}
	if assignedUserID != "" {
		query += ` JOIN service_order_assignments a ON a.service_order_id = o.id AND a.status = 'active'`
	}
	query += ` WHERE o.business_id = $1`
	if branchID != "" {
		args = append(args, branchID)
		query += ` AND o.branch_id = $` + strconv.Itoa(len(args))
	}
	if assignedUserID != "" {
		args = append(args, assignedUserID)
		query += ` AND a.assigned_user_id = $` + strconv.Itoa(len(args))
	}
	query += ` GROUP BY o.status`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return domain.ServiceOrderSummary{}, err
	}
	defer rows.Close()

	summary := domain.ServiceOrderSummary{}
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return domain.ServiceOrderSummary{}, err
		}
		summary.Total += count
		switch status {
		case "open":
			summary.Open = count
		case "in_progress":
			summary.InProgress = count
		case "completed":
			summary.Completed = count
		case "cancelled":
			summary.Cancelled = count
		}
	}
	return summary, rows.Err()
}

func (r *ServiceOrderRepository) FindByID(ctx context.Context, serviceOrderID string) (domain.ServiceOrder, error) {
	var order domain.ServiceOrder
	err := r.db.QueryRow(ctx, `
		SELECT id, business_id, branch_id, service_definition_id, title, description, status, priority
		FROM service_orders
		WHERE id = $1
	`, serviceOrderID).Scan(&order.ID, &order.BusinessID, &order.BranchID, &order.ServiceDefinitionID, &order.Title, &order.Description, &order.Status, &order.Priority)
	return order, err
}

func (r *ServiceOrderRepository) UpdateStatus(ctx context.Context, order domain.ServiceOrder, previousStatus string, changedByUserID string, requestID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		UPDATE service_orders
		SET status = $2, updated_at = now()
		WHERE id = $1
	`, order.ID, order.Status); err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO service_order_status_history (id, service_order_id, business_id, branch_id, from_status, to_status, changed_by_user_id, request_id)
		VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7, '')::uuid, $8)
	`, uuid.New(), order.ID, order.BusinessID, order.BranchID, previousStatus, order.Status, changedByUserID, requestID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *ServiceOrderRepository) Assign(ctx context.Context, assignment domain.ServiceOrderAssignment) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		UPDATE service_order_assignments
		SET status = 'replaced', updated_at = now()
		WHERE service_order_id = $1 AND status = 'active'
	`, assignment.ServiceOrderID); err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO service_order_assignments (id, service_order_id, business_id, branch_id, assigned_user_id, assigned_by_user_id, status, request_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, assignment.ID, assignment.ServiceOrderID, assignment.BusinessID, assignment.BranchID, assignment.AssignedUserID, assignment.AssignedByUserID, assignment.Status, assignment.RequestID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *ServiceOrderRepository) ListAssignedToUser(ctx context.Context, businessID string, assignedUserID string, branchID string) ([]domain.ServiceOrder, error) {
	query := `
		SELECT o.id, o.business_id, o.branch_id, o.service_definition_id, o.title, o.description, o.status, o.priority
		FROM service_orders o
		JOIN service_order_assignments a ON a.service_order_id = o.id AND a.status = 'active'
		WHERE o.business_id = $1 AND a.assigned_user_id = $2`
	args := []any{businessID, assignedUserID}
	if branchID != "" {
		query += ` AND o.branch_id = $3`
		args = append(args, branchID)
	}
	query += ` ORDER BY o.created_at DESC`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []domain.ServiceOrder{}
	for rows.Next() {
		var order domain.ServiceOrder
		if err := rows.Scan(&order.ID, &order.BusinessID, &order.BranchID, &order.ServiceDefinitionID, &order.Title, &order.Description, &order.Status, &order.Priority); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, rows.Err()
}

func (r *ServiceOrderRepository) ListAssignmentsByOrder(ctx context.Context, serviceOrderID string) ([]domain.ServiceOrderAssignment, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, service_order_id, business_id, branch_id, assigned_user_id, assigned_by_user_id, status, request_id
		FROM service_order_assignments
		WHERE service_order_id = $1 AND status = 'active'
		ORDER BY created_at DESC
	`, serviceOrderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assignments := []domain.ServiceOrderAssignment{}
	for rows.Next() {
		var assignment domain.ServiceOrderAssignment
		if err := rows.Scan(&assignment.ID, &assignment.ServiceOrderID, &assignment.BusinessID, &assignment.BranchID, &assignment.AssignedUserID, &assignment.AssignedByUserID, &assignment.Status, &assignment.RequestID); err != nil {
			return nil, err
		}
		assignments = append(assignments, assignment)
	}
	return assignments, rows.Err()
}
