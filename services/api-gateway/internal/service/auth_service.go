package service

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/isapr/mini-erp/services/api-gateway/internal/domain"
	"github.com/isapr/mini-erp/services/api-gateway/internal/ports"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrValidation         = errors.New("validation error")
)

type account struct {
	userID       string
	businessID   string
	fullName     string
	email        string
	password     string
	role         domain.Role
	permissions  []string
	refreshToken string
	accessToken  string
}

type AuthService struct {
	mu           sync.RWMutex
	accounts     map[string]*account
	tokens       map[string]string
	refresh      map[string]string
	identity     ports.IdentityClient
	organization ports.OrganizationClient
}

func NewAuthService() *AuthService {
	service := &AuthService{
		accounts: make(map[string]*account),
		tokens:   make(map[string]string),
		refresh:  make(map[string]string),
	}
	service.seedPlatformAdmin()
	return service
}

func (s *AuthService) seedPlatformAdmin() {
	acct := &account{
		userID:      uuid.NewString(),
		fullName:    "Platform Admin",
		email:       "platform.admin@example.com",
		password:    "admin123",
		role:        domain.RolePlatformAdmin,
		permissions: []string{"platform.businesses.read", "platform.businesses.update"},
	}
	s.accounts[strings.ToLower(acct.email)] = acct
}

func NewAuthServiceWithClients(identity ports.IdentityClient, organization ports.OrganizationClient) *AuthService {
	service := NewAuthService()
	service.identity = identity
	service.organization = organization
	return service
}

func (s *AuthService) Signup(ctx context.Context, req domain.SignupRequest, requestID string) (domain.AuthSession, error) {
	if req.BusinessName == "" || req.BusinessTimezone == "" || req.AdminFullName == "" || req.AdminEmail == "" || req.AdminPassword == "" {
		return domain.AuthSession{}, ErrValidation
	}
	if s.identity != nil && s.organization != nil {
		return s.signupWithGRPC(ctx, req, requestID)
	}

	key := strings.ToLower(strings.TrimSpace(req.AdminEmail))

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.accounts[key]; exists {
		return domain.AuthSession{}, ErrValidation
	}

	acct := &account{
		userID:       uuid.NewString(),
		businessID:   uuid.NewString(),
		fullName:     req.AdminFullName,
		email:        key,
		password:     req.AdminPassword,
		role:         domain.RoleBusinessAdmin,
		permissions:  []string{"business.manage", "branches.manage", "users.manage", "workflows.manage", "service-orders.manage", "resources.manage", "reports.read"},
		accessToken:  uuid.NewString(),
		refreshToken: uuid.NewString(),
	}

	s.accounts[key] = acct
	s.tokens[acct.accessToken] = key
	s.refresh[acct.refreshToken] = key

	return sessionFromAccount(acct, requestID), nil
}

func (s *AuthService) signupWithGRPC(ctx context.Context, req domain.SignupRequest, requestID string) (domain.AuthSession, error) {
	business, err := s.organization.CreateBusiness(ctx, ports.CreateBusinessRequest{Name: req.BusinessName, Timezone: req.BusinessTimezone})
	if err != nil {
		return domain.AuthSession{}, err
	}

	admin, err := s.identity.SignupTenantAdmin(ctx, ports.SignupTenantAdminRequest{
		BusinessID:    business.BusinessID,
		AdminEmail:    req.AdminEmail,
		AdminPassword: req.AdminPassword,
		AdminFullName: req.AdminFullName,
		RequestID:     requestID,
	})
	if err != nil {
		return domain.AuthSession{}, err
	}

	return domain.AuthSession{
		AccessToken:  admin.AccessToken,
		RefreshToken: admin.RefreshToken,
		User: domain.AuthContext{
			UserID:            admin.UserID,
			BusinessID:        admin.BusinessID,
			Role:              domain.Role(admin.Role),
			Permissions:       admin.Permissions,
			AssignedBranchIDs: admin.AssignedBranchIDs,
			RequestID:         admin.RequestID,
		},
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req domain.LoginRequest, requestID string) (domain.AuthSession, error) {
	if s.identity != nil {
		resp, err := s.identity.Login(ctx, ports.LoginRequest{Email: req.Email, Password: req.Password, RequestID: requestID})
		if err != nil {
			return domain.AuthSession{}, err
		}
		return sessionFromIdentity(resp), nil
	}

	key := strings.ToLower(strings.TrimSpace(req.Email))

	s.mu.Lock()
	defer s.mu.Unlock()

	acct, ok := s.accounts[key]
	if !ok || acct.password != req.Password {
		return domain.AuthSession{}, ErrInvalidCredentials
	}

	acct.accessToken = uuid.NewString()
	acct.refreshToken = uuid.NewString()
	s.tokens[acct.accessToken] = key
	s.refresh[acct.refreshToken] = key

	return sessionFromAccount(acct, requestID), nil
}

func (s *AuthService) Refresh(_ context.Context, refreshToken string, requestID string) (domain.AuthSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key, ok := s.refresh[refreshToken]
	if !ok {
		return domain.AuthSession{}, ErrInvalidToken
	}

	acct := s.accounts[key]
	acct.accessToken = uuid.NewString()
	acct.refreshToken = uuid.NewString()
	s.tokens[acct.accessToken] = key
	s.refresh[acct.refreshToken] = key

	return sessionFromAccount(acct, requestID), nil
}

func (s *AuthService) GetMe(ctx context.Context, accessToken string, requestID string) (domain.AuthContext, error) {
	if s.identity != nil {
		resp, err := s.identity.GetUserAccessContext(ctx, ports.GetUserAccessContextRequest{AccessToken: accessToken, RequestID: requestID})
		if err != nil {
			return domain.AuthContext{}, err
		}
		authContext := domain.AuthContext{UserID: resp.UserID, BusinessID: resp.BusinessID, Role: domain.Role(resp.Role), Permissions: resp.Permissions, AssignedBranchIDs: resp.AssignedBranchIDs, RequestID: resp.RequestID}
		if s.organization != nil && authContext.BusinessID != "" && (authContext.Role == domain.RoleManager || authContext.Role == domain.RoleStaff) {
			assigned, err := s.organization.ListAssignedBranches(ctx, ports.ListAssignedBranchesRequest{UserID: authContext.UserID, BusinessID: authContext.BusinessID})
			if err != nil {
				return domain.AuthContext{}, err
			}
			authContext.AssignedBranchIDs = assigned.BranchIDs
		}
		return authContext, nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	key, ok := s.tokens[accessToken]
	if !ok {
		return domain.AuthContext{}, ErrInvalidToken
	}

	return contextFromAccount(s.accounts[key], requestID), nil
}

func sessionFromIdentity(resp ports.SignupTenantAdminResponse) domain.AuthSession {
	return domain.AuthSession{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		User: domain.AuthContext{
			UserID:            resp.UserID,
			BusinessID:        resp.BusinessID,
			Role:              domain.Role(resp.Role),
			Permissions:       resp.Permissions,
			AssignedBranchIDs: resp.AssignedBranchIDs,
			RequestID:         resp.RequestID,
		},
	}
}

func (s *AuthService) Logout(_ context.Context, accessToken string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.tokens, accessToken)
	return nil
}

func sessionFromAccount(acct *account, requestID string) domain.AuthSession {
	return domain.AuthSession{
		AccessToken:  acct.accessToken,
		RefreshToken: acct.refreshToken,
		User:         contextFromAccount(acct, requestID),
	}
}

func contextFromAccount(acct *account, requestID string) domain.AuthContext {
	ctx := domain.AuthContext{
		UserID:      acct.userID,
		Role:        acct.role,
		Permissions: acct.permissions,
		RequestID:   requestID,
	}
	if acct.businessID != "" {
		ctx.BusinessID = acct.businessID
		ctx.AssignedBranchIDs = []string{}
	}
	return ctx
}
