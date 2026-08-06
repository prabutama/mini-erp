package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/isapr/mini-erp/services/identity/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	FindByID(ctx context.Context, userID uuid.UUID) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindAccessContext(ctx context.Context, userID uuid.UUID, requestID string) (domain.AuthContext, error)
	ListByBusiness(ctx context.Context, businessID uuid.UUID) ([]domain.User, error)
	Update(ctx context.Context, user domain.User) error
	AssignBusinessRole(ctx context.Context, userID uuid.UUID, roleName domain.RoleName, businessID uuid.UUID) error
	CreateRefreshToken(ctx context.Context, tokenID uuid.UUID, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
}
