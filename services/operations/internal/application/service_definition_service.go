package application

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/isapr/mini-erp/services/operations/internal/domain"
	"github.com/isapr/mini-erp/services/operations/internal/ports"
)

var ErrValidation = errors.New("validation error")

type CreateServiceDefinitionInput struct {
	BusinessID  string
	Name        string
	Description string
}

type ServiceDefinitionService struct {
	services ports.ServiceDefinitionRepository
}

func NewServiceDefinitionService(services ports.ServiceDefinitionRepository) *ServiceDefinitionService {
	return &ServiceDefinitionService{services: services}
}

func (s *ServiceDefinitionService) Create(ctx context.Context, input CreateServiceDefinitionInput) (domain.ServiceDefinition, error) {
	businessID, err := uuid.Parse(input.BusinessID)
	if err != nil {
		return domain.ServiceDefinition{}, ErrValidation
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return domain.ServiceDefinition{}, ErrValidation
	}
	service := domain.ServiceDefinition{ID: uuid.New(), BusinessID: businessID, Name: name, Code: slug(name), Description: strings.TrimSpace(input.Description), Status: "active"}
	if err := s.services.Create(ctx, service); err != nil {
		return domain.ServiceDefinition{}, err
	}
	return service, nil
}

func (s *ServiceDefinitionService) List(ctx context.Context, businessID string) ([]domain.ServiceDefinition, error) {
	if _, err := uuid.Parse(businessID); err != nil {
		return nil, ErrValidation
	}
	return s.services.ListByBusiness(ctx, businessID)
}

func (s *ServiceDefinitionService) Get(ctx context.Context, serviceDefinitionID string) (domain.ServiceDefinition, error) {
	if _, err := uuid.Parse(serviceDefinitionID); err != nil {
		return domain.ServiceDefinition{}, ErrValidation
	}
	return s.services.FindByID(ctx, serviceDefinitionID)
}

func slug(value string) string {
	clean := regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(strings.ToLower(strings.TrimSpace(value)), "-")
	return strings.Trim(clean, "-")
}
