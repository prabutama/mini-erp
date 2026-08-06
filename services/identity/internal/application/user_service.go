package application

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/isapr/mini-erp/services/identity/internal/domain"
	"github.com/isapr/mini-erp/services/identity/internal/ports"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserInput struct {
	BusinessID uuid.UUID
	Email      string
	Password   string
	FullName   string
	Role       domain.RoleName
}

type UpdateUserInput struct {
	UserID   uuid.UUID
	FullName string
	Status   string
}

type UserService struct {
	users ports.UserRepository
}

func NewUserService(users ports.UserRepository) *UserService {
	return &UserService{users: users}
}

func (s *UserService) CreateBusinessUser(ctx context.Context, input CreateUserInput) (domain.User, error) {
	if input.BusinessID == uuid.Nil || input.Email == "" || input.Password == "" || input.FullName == "" || !isBusinessRole(input.Role) {
		return domain.User{}, ErrValidation
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
	}

	user := domain.User{ID: uuid.New(), Email: strings.ToLower(strings.TrimSpace(input.Email)), PasswordHash: string(passwordHash), FullName: strings.TrimSpace(input.FullName), Status: "active"}
	if err := s.users.Create(ctx, user); err != nil {
		return domain.User{}, err
	}
	if err := s.users.AssignBusinessRole(ctx, user.ID, input.Role, input.BusinessID); err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (s *UserService) ListBusinessUsers(ctx context.Context, businessID uuid.UUID) ([]domain.User, error) {
	if businessID == uuid.Nil {
		return nil, ErrValidation
	}
	return s.users.ListByBusiness(ctx, businessID)
}

func (s *UserService) GetUser(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	return s.users.FindByID(ctx, userID)
}

func (s *UserService) UpdateUser(ctx context.Context, input UpdateUserInput) (domain.User, error) {
	user, err := s.users.FindByID(ctx, input.UserID)
	if err != nil {
		return domain.User{}, err
	}
	if input.FullName != "" {
		user.FullName = strings.TrimSpace(input.FullName)
	}
	if input.Status != "" {
		user.Status = strings.TrimSpace(input.Status)
	}
	if err := s.users.Update(ctx, user); err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (s *UserService) AssignBusinessRole(ctx context.Context, userID uuid.UUID, role domain.RoleName, businessID uuid.UUID) error {
	if userID == uuid.Nil || businessID == uuid.Nil || !isBusinessRole(role) {
		return ErrValidation
	}
	return s.users.AssignBusinessRole(ctx, userID, role, businessID)
}

func isBusinessRole(role domain.RoleName) bool {
	return role == domain.RoleBusinessAdmin || role == domain.RoleManager || role == domain.RoleStaff
}
