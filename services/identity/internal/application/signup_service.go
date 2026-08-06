package application

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/isapr/mini-erp/services/identity/internal/domain"
	"github.com/isapr/mini-erp/services/identity/internal/ports"
	"golang.org/x/crypto/bcrypt"
)

var ErrValidation = errors.New("validation error")
var ErrUserAlreadyExists = domain.ErrUserAlreadyExists

type SignupTenantAdminInput struct {
	BusinessID uuid.UUID
	Email      string
	Password   string
	FullName   string
}

type SignupTenantAdminOutput struct {
	UserID     uuid.UUID
	BusinessID uuid.UUID
	Role       domain.RoleName
}

type SignupService struct {
	users ports.UserRepository
	auth  *AuthService
}

func NewSignupService(users ports.UserRepository) *SignupService {
	return &SignupService{users: users, auth: NewAuthService(users, "")}
}

func NewSignupServiceWithAuth(users ports.UserRepository, auth *AuthService) *SignupService {
	return &SignupService{users: users, auth: auth}
}

func (s *SignupService) SignupTenantAdmin(ctx context.Context, input SignupTenantAdminInput) (SignupTenantAdminOutput, error) {
	if input.BusinessID == uuid.Nil || input.Email == "" || input.Password == "" || input.FullName == "" {
		return SignupTenantAdminOutput{}, ErrValidation
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return SignupTenantAdminOutput{}, err
	}

	user := domain.User{
		ID:           uuid.New(),
		Email:        strings.ToLower(strings.TrimSpace(input.Email)),
		PasswordHash: string(passwordHash),
		FullName:     strings.TrimSpace(input.FullName),
		Status:       "active",
	}

	if err := s.users.Create(ctx, user); err != nil {
		return SignupTenantAdminOutput{}, err
	}

	if err := s.users.AssignBusinessRole(ctx, user.ID, domain.RoleBusinessAdmin, input.BusinessID); err != nil {
		return SignupTenantAdminOutput{}, err
	}

	return SignupTenantAdminOutput{UserID: user.ID, BusinessID: input.BusinessID, Role: domain.RoleBusinessAdmin}, nil
}

func (s *SignupService) SignupTenantAdminSession(ctx context.Context, input SignupTenantAdminInput, requestID string) (domain.AuthSession, error) {
	output, err := s.SignupTenantAdmin(ctx, input)
	if err != nil {
		return domain.AuthSession{}, err
	}
	return s.auth.SessionForUser(ctx, output.UserID, requestID)
}
