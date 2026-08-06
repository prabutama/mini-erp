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

func (r *PlacementRepository) ListAssignedBranches(ctx context.Context, userID string, businessID string) ([]domain.EmployeePlacement, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, business_id, branch_id, position, employment_type, status, start_date
		FROM employee_placements
		WHERE user_id = $1 AND business_id = $2 AND status = 'active' AND end_date IS NULL
		ORDER BY start_date DESC, created_at DESC
	`, userID, businessID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	placements := []domain.EmployeePlacement{}
	for rows.Next() {
		var placement domain.EmployeePlacement
		if err := rows.Scan(&placement.ID, &placement.UserID, &placement.BusinessID, &placement.BranchID, &placement.Position, &placement.EmploymentType, &placement.Status, &placement.StartDate); err != nil {
			return nil, err
		}
		placements = append(placements, placement)
	}
	return placements, rows.Err()
}
