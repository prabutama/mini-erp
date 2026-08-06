package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/isapr/mini-erp/services/identity/internal/domain"
	"github.com/isapr/mini-erp/services/identity/internal/ports"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrInvalidToken = errors.New("invalid token")

type AuthService struct {
	users     ports.UserRepository
	jwtSecret []byte
}

func NewAuthService(users ports.UserRepository, jwtSecret string) *AuthService {
	if jwtSecret == "" {
		jwtSecret = "dev-secret"
	}
	return &AuthService{users: users, jwtSecret: []byte(jwtSecret)}
}

type LoginInput struct {
	Email     string
	Password  string
	RequestID string
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (domain.AuthSession, error) {
	if input.Email == "" || input.Password == "" {
		return domain.AuthSession{}, ErrValidation
	}

	user, err := s.users.FindByEmail(ctx, strings.ToLower(strings.TrimSpace(input.Email)))
	if err != nil {
		return domain.AuthSession{}, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return domain.AuthSession{}, ErrInvalidCredentials
	}

	return s.sessionForUser(ctx, user.ID, input.RequestID)
}

func (s *AuthService) SessionForUser(ctx context.Context, userID uuid.UUID, requestID string) (domain.AuthSession, error) {
	return s.sessionForUser(ctx, userID, requestID)
}

func (s *AuthService) GetUserAccessContext(ctx context.Context, accessToken string, requestID string) (domain.AuthContext, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return domain.AuthContext{}, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return domain.AuthContext{}, ErrInvalidToken
	}

	userIDValue, ok := claims["user_id"].(string)
	if !ok {
		return domain.AuthContext{}, ErrInvalidToken
	}
	userID, err := uuid.Parse(userIDValue)
	if err != nil {
		return domain.AuthContext{}, ErrInvalidToken
	}

	return s.users.FindAccessContext(ctx, userID, requestID)
}

func (s *AuthService) sessionForUser(ctx context.Context, userID uuid.UUID, requestID string) (domain.AuthSession, error) {
	authContext, err := s.users.FindAccessContext(ctx, userID, requestID)
	if err != nil {
		return domain.AuthSession{}, err
	}

	accessToken, err := s.accessToken(authContext)
	if err != nil {
		return domain.AuthSession{}, err
	}

	refreshToken := uuid.NewString()
	if err := s.users.CreateRefreshToken(ctx, uuid.New(), userID, hashToken(refreshToken), time.Now().Add(30*24*time.Hour)); err != nil {
		return domain.AuthSession{}, err
	}

	return domain.AuthSession{AccessToken: accessToken, RefreshToken: refreshToken, Context: authContext}, nil
}

func (s *AuthService) accessToken(authContext domain.AuthContext) (string, error) {
	claims := jwt.MapClaims{
		"user_id": authContext.UserID.String(),
		"role":    string(authContext.Role),
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	if authContext.BusinessID != uuid.Nil {
		claims["business_id"] = authContext.BusinessID.String()
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
