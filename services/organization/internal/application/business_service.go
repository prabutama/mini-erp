package application

import (
	"context"
	"errors"
	"regexp"
	"strings"

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

func businessCode(name string) string {
	code := strings.ToLower(strings.TrimSpace(name))
	code = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(code, "-")
	code = strings.Trim(code, "-")
	if code == "" {
		return uuid.NewString()
	}
	return code
}
