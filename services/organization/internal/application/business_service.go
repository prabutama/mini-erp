package application

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/isapr/mini-erp/services/organization/internal/domain"
	"github.com/isapr/mini-erp/services/organization/internal/ports"
)

var ErrValidation = errors.New("validation error")
var ErrBusinessAlreadyExists = domain.ErrBusinessAlreadyExists

type CreateBusinessInput struct {
	Name     string
	Timezone string
}

type UpdateBusinessInput struct {
	BusinessID string
	Name       string
	Timezone   string
}

type UpdatePlatformBusinessInput struct {
	BusinessID     string
	Status         string
	Plan           string
	PlatformNotes  string
	SuspendedAtRFC string
}

type BusinessService struct {
	businesses ports.BusinessRepository
}

func NewBusinessService(businesses ports.BusinessRepository) *BusinessService {
	return &BusinessService{businesses: businesses}
}

func (s *BusinessService) CreateBusiness(ctx context.Context, input CreateBusinessInput) (domain.Business, error) {
	if input.Name == "" || input.Timezone == "" {
		return domain.Business{}, ErrValidation
	}

	business := domain.Business{
		ID:       uuid.New(),
		Name:     strings.TrimSpace(input.Name),
		Code:     businessCode(input.Name),
		Status:   "active",
		Plan:     "free",
		Timezone: strings.TrimSpace(input.Timezone),
	}

	if err := s.businesses.Create(ctx, business); err != nil {
		return domain.Business{}, err
	}

	return business, nil
}

func (s *BusinessService) GetBusiness(ctx context.Context, businessID string) (domain.Business, error) {
	if businessID == "" {
		return domain.Business{}, ErrValidation
	}
	return s.businesses.FindByID(ctx, businessID)
}

func (s *BusinessService) ListPlatformBusinesses(ctx context.Context) ([]domain.Business, error) {
	return s.businesses.List(ctx)
}

func (s *BusinessService) UpdateBusiness(ctx context.Context, input UpdateBusinessInput) (domain.Business, error) {
	if input.BusinessID == "" {
		return domain.Business{}, ErrValidation
	}
	business, err := s.businesses.FindByID(ctx, input.BusinessID)
	if err != nil {
		return domain.Business{}, err
	}
	if strings.TrimSpace(input.Name) != "" {
		business.Name = strings.TrimSpace(input.Name)
	}
	if strings.TrimSpace(input.Timezone) != "" {
		business.Timezone = strings.TrimSpace(input.Timezone)
	}
	if err := s.businesses.Update(ctx, business); err != nil {
		return domain.Business{}, err
	}
	return s.businesses.FindByID(ctx, input.BusinessID)
}

func (s *BusinessService) UpdatePlatformBusiness(ctx context.Context, input UpdatePlatformBusinessInput) (domain.Business, error) {
	if input.BusinessID == "" {
		return domain.Business{}, ErrValidation
	}
	business, err := s.businesses.FindByID(ctx, input.BusinessID)
	if err != nil {
		return domain.Business{}, err
	}
	if strings.TrimSpace(input.Status) != "" {
		business.Status = strings.TrimSpace(input.Status)
	}
	if strings.TrimSpace(input.Plan) != "" {
		business.Plan = strings.TrimSpace(input.Plan)
	}
	business.PlatformNotes = strings.TrimSpace(input.PlatformNotes)
	if strings.TrimSpace(input.SuspendedAtRFC) != "" {
		parsed, err := time.Parse(time.RFC3339, input.SuspendedAtRFC)
		if err != nil {
			return domain.Business{}, ErrValidation
		}
		business.SuspendedAt = &parsed
	}
	if err := s.businesses.UpdatePlatform(ctx, business); err != nil {
		return domain.Business{}, err
	}
	return s.businesses.FindByID(ctx, input.BusinessID)
}

func businessCode(name string) string {
	code := strings.ToLower(strings.TrimSpace(name))
	code = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(code, "-")
	code = strings.Trim(code, "-")
	if code == "" {
		return uuid.NewString()
	}
	return code
}
