package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/isapr/mini-erp/services/organization/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

type BusinessRepository struct {
	db *pgxpool.Pool
}

func NewBusinessRepository(db *pgxpool.Pool) *BusinessRepository {
	return &BusinessRepository{db: db}
}

func (r *BusinessRepository) Create(ctx context.Context, business domain.Business) error {
	_, err := r.db.Exec(ctx, `
        INSERT INTO businesses (id, name, code, status, plan, timezone)
        VALUES ($1, $2, $3, $4, $5, $6)
	`, business.ID, business.Name, strings.ToLower(business.Code), business.Status, business.Plan, business.Timezone)
	if isUniqueViolation(err) {
		return domain.ErrBusinessAlreadyExists
	}
	return err
}

func (r *BusinessRepository) FindByCode(ctx context.Context, code string) (domain.Business, error) {
	var business domain.Business
	err := r.db.QueryRow(ctx, `
		SELECT id, name, code, status, plan, timezone, created_at
		FROM businesses
		WHERE code = $1
	`, strings.ToLower(code)).Scan(&business.ID, &business.Name, &business.Code, &business.Status, &business.Plan, &business.Timezone, &business.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Business{}, ErrNotFound
	}
	return business, err
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
