package postgres

import (
	"context"

	"github.com/isapr/mini-erp/services/organization/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlacementRepository struct {
	db *pgxpool.Pool
}

func NewPlacementRepository(db *pgxpool.Pool) *PlacementRepository {
	return &PlacementRepository{db: db}
}

func (r *PlacementRepository) Create(ctx context.Context, placement domain.EmployeePlacement) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO employee_placements (id, user_id, business_id, branch_id, position, employment_type, status, start_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, placement.ID, placement.UserID, placement.BusinessID, placement.BranchID, placement.Position, placement.EmploymentType, placement.Status, placement.StartDate)
	return err
}
