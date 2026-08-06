package postgres

import (
	"context"
	"errors"

	"github.com/isapr/mini-erp/services/organization/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BranchRepository struct {
	db *pgxpool.Pool
}

func NewBranchRepository(db *pgxpool.Pool) *BranchRepository {
	return &BranchRepository{db: db}
}

func (r *BranchRepository) Create(ctx context.Context, branch domain.Branch) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO branches (id, business_id, name, code, address, phone, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, branch.ID, branch.BusinessID, branch.Name, branch.Code, branch.Address, branch.Phone, branch.Status)
	return err
}

func (r *BranchRepository) FindByID(ctx context.Context, branchID string) (domain.Branch, error) {
	var branch domain.Branch
	err := r.db.QueryRow(ctx, `
		SELECT id, business_id, name, code, COALESCE(address, ''), COALESCE(phone, ''), status, created_at
		FROM branches
		WHERE id = $1
	`, branchID).Scan(&branch.ID, &branch.BusinessID, &branch.Name, &branch.Code, &branch.Address, &branch.Phone, &branch.Status, &branch.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Branch{}, ErrNotFound
	}
	return branch, err
}

func (r *BranchRepository) ListByBusiness(ctx context.Context, businessID string) ([]domain.Branch, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, business_id, name, code, COALESCE(address, ''), COALESCE(phone, ''), status, created_at
		FROM branches
		WHERE business_id = $1
		ORDER BY created_at DESC
	`, businessID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	branches := []domain.Branch{}
	for rows.Next() {
		var branch domain.Branch
		if err := rows.Scan(&branch.ID, &branch.BusinessID, &branch.Name, &branch.Code, &branch.Address, &branch.Phone, &branch.Status, &branch.CreatedAt); err != nil {
			return nil, err
		}
		branches = append(branches, branch)
	}
	return branches, rows.Err()
}

func (r *BranchRepository) Update(ctx context.Context, branch domain.Branch) error {
	_, err := r.db.Exec(ctx, `
		UPDATE branches
		SET name = $2, address = $3, phone = $4, status = $5, updated_at = now()
		WHERE id = $1
	`, branch.ID, branch.Name, branch.Address, branch.Phone, branch.Status)
	return err
}
