package postgres

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/isapr/mini-erp/services/identity/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) error {
	_, err := r.db.Exec(ctx, `
        INSERT INTO users (id, email, password_hash, full_name, status)
        VALUES ($1, $2, $3, $4, $5)
	`, user.ID, strings.ToLower(user.Email), user.PasswordHash, user.FullName, user.Status)
	if isUniqueViolation(err) {
		return domain.ErrUserAlreadyExists
	}
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(ctx, `
		SELECT id, email, password_hash, full_name, status, created_at
		FROM users
		WHERE email = $1
	`, strings.ToLower(email)).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Status, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, ErrNotFound
	}
	return user, err
}

func (r *UserRepository) FindByID(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(ctx, `
		SELECT id, email, password_hash, full_name, status, created_at
		FROM users
		WHERE id = $1
	`, userID).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Status, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, ErrNotFound
	}
	return user, err
}

func (r *UserRepository) ListByBusiness(ctx context.Context, businessID uuid.UUID) ([]domain.User, error) {
	rows, err := r.db.Query(ctx, `
		SELECT users.id, users.email, users.password_hash, users.full_name, users.status, users.created_at
		FROM users
		JOIN user_roles ON user_roles.user_id = users.id
		WHERE user_roles.business_id = $1
		ORDER BY users.created_at DESC
	`, businessID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []domain.User{}
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Status, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r *UserRepository) Update(ctx context.Context, user domain.User) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users
		SET full_name = $2, status = $3, updated_at = now()
		WHERE id = $1
	`, user.ID, user.FullName, user.Status)
	return err
}

func (r *UserRepository) FindAccessContext(ctx context.Context, userID uuid.UUID, requestID string) (domain.AuthContext, error) {
	var authContext domain.AuthContext
	authContext.UserID = userID
	authContext.RequestID = requestID

	err := r.db.QueryRow(ctx, `
		SELECT roles.name, COALESCE(user_roles.business_id, '00000000-0000-0000-0000-000000000000'::uuid)
		FROM user_roles
		JOIN roles ON roles.id = user_roles.role_id
		WHERE user_roles.user_id = $1
		ORDER BY CASE roles.scope WHEN 'platform' THEN 0 ELSE 1 END
		LIMIT 1
	`, userID).Scan(&authContext.Role, &authContext.BusinessID)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.AuthContext{}, ErrNotFound
	}
	if err != nil {
		return domain.AuthContext{}, err
	}

	authContext.Permissions = permissionsForRole(authContext.Role)
	if authContext.Role == domain.RolePlatformAdmin {
		authContext.BusinessID = uuid.Nil
	}
	if authContext.BusinessID != uuid.Nil {
		authContext.AssignedBranchIDs = []string{}
	}

	return authContext, nil
}

func (r *UserRepository) AssignBusinessRole(ctx context.Context, userID uuid.UUID, roleName domain.RoleName, businessID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO user_roles (user_id, role_id, business_id)
		SELECT $1, roles.id, $3
		FROM roles
		WHERE roles.name = $2 AND roles.scope = 'business'
	`, userID, string(roleName), businessID)
	return err
}

func (r *UserRepository) CreateRefreshToken(ctx context.Context, tokenID uuid.UUID, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at)
		VALUES ($1, $2, $3, $4)
	`, tokenID, userID, tokenHash, expiresAt)
	return err
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func permissionsForRole(role domain.RoleName) []string {
	switch role {
	case domain.RolePlatformAdmin:
		return []string{"platform.businesses.read", "platform.businesses.update"}
	case domain.RoleBusinessAdmin:
		return []string{"business.manage", "branches.manage", "users.manage", "workflows.manage", "service-orders.manage", "resources.manage", "reports.read"}
	case domain.RoleManager:
		return []string{"service-orders.manage", "resources.manage"}
	case domain.RoleStaff:
		return []string{"service-orders.update", "resources.use"}
	default:
		return []string{}
	}
}
