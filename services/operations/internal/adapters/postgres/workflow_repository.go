package postgres

import (
	"context"

	"github.com/isapr/mini-erp/services/operations/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkflowRepository struct{ db *pgxpool.Pool }

func NewWorkflowRepository(db *pgxpool.Pool) *WorkflowRepository { return &WorkflowRepository{db: db} }

func (r *WorkflowRepository) Create(ctx context.Context, workflow domain.Workflow) error {
	_, err := r.db.Exec(ctx, `INSERT INTO workflows (id, business_id, name, description, status) VALUES ($1,$2,$3,$4,$5)`, workflow.ID, workflow.BusinessID, workflow.Name, workflow.Description, workflow.Status)
	return err
}

func (r *WorkflowRepository) ListByBusiness(ctx context.Context, businessID string) ([]domain.Workflow, error) {
	rows, err := r.db.Query(ctx, `SELECT id, business_id, name, description, status FROM workflows WHERE business_id = $1 ORDER BY name`, businessID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []domain.Workflow{}
	for rows.Next() {
		var item domain.Workflow
		if err := rows.Scan(&item.ID, &item.BusinessID, &item.Name, &item.Description, &item.Status); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *WorkflowRepository) FindByID(ctx context.Context, workflowID string) (domain.Workflow, error) {
	var item domain.Workflow
	err := r.db.QueryRow(ctx, `SELECT id, business_id, name, description, status FROM workflows WHERE id = $1`, workflowID).Scan(&item.ID, &item.BusinessID, &item.Name, &item.Description, &item.Status)
	return item, err
}

func (r *WorkflowRepository) Update(ctx context.Context, workflow domain.Workflow) error {
	_, err := r.db.Exec(ctx, `UPDATE workflows SET name = $2, description = $3, status = $4, updated_at = now() WHERE id = $1`, workflow.ID, workflow.Name, workflow.Description, workflow.Status)
	return err
}

func (r *WorkflowRepository) CreateStatus(ctx context.Context, status domain.WorkflowStatus) error {
	_, err := r.db.Exec(ctx, `INSERT INTO workflow_statuses (id, workflow_id, business_id, code, name, category, sort_order, is_initial, is_terminal) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`, status.ID, status.WorkflowID, status.BusinessID, status.Code, status.Name, status.Category, status.SortOrder, status.IsInitial, status.IsTerminal)
	return err
}

func (r *WorkflowRepository) ListStatuses(ctx context.Context, workflowID string) ([]domain.WorkflowStatus, error) {
	rows, err := r.db.Query(ctx, `SELECT id, workflow_id, business_id, code, name, category, sort_order, is_initial, is_terminal FROM workflow_statuses WHERE workflow_id = $1 ORDER BY sort_order, code`, workflowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []domain.WorkflowStatus{}
	for rows.Next() {
		var item domain.WorkflowStatus
		if err := rows.Scan(&item.ID, &item.WorkflowID, &item.BusinessID, &item.Code, &item.Name, &item.Category, &item.SortOrder, &item.IsInitial, &item.IsTerminal); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *WorkflowRepository) CreateTransition(ctx context.Context, transition domain.WorkflowTransition) error {
	_, err := r.db.Exec(ctx, `INSERT INTO workflow_transitions (id, workflow_id, business_id, from_status_code, to_status_code) VALUES ($1,$2,$3,$4,$5)`, transition.ID, transition.WorkflowID, transition.BusinessID, transition.FromStatusCode, transition.ToStatusCode)
	return err
}

func (r *WorkflowRepository) ListTransitions(ctx context.Context, workflowID string) ([]domain.WorkflowTransition, error) {
	rows, err := r.db.Query(ctx, `SELECT id, workflow_id, business_id, from_status_code, to_status_code FROM workflow_transitions WHERE workflow_id = $1 ORDER BY from_status_code, to_status_code`, workflowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []domain.WorkflowTransition{}
	for rows.Next() {
		var item domain.WorkflowTransition
		if err := rows.Scan(&item.ID, &item.WorkflowID, &item.BusinessID, &item.FromStatusCode, &item.ToStatusCode); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
