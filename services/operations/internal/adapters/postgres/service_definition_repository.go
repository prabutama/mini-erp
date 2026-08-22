package postgres

import (
	"context"

	"github.com/isapr/mini-erp/services/operations/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceDefinitionRepository struct {
	db *pgxpool.Pool
}

func NewServiceDefinitionRepository(db *pgxpool.Pool) *ServiceDefinitionRepository {
	return &ServiceDefinitionRepository{db: db}
}

func (r *ServiceDefinitionRepository) Create(ctx context.Context, service domain.ServiceDefinition) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO service_definitions (id, business_id, name, code, description, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, service.ID, service.BusinessID, service.Name, service.Code, service.Description, service.Status)
	return err
}

func (r *ServiceDefinitionRepository) ListByBusiness(ctx context.Context, businessID string) ([]domain.ServiceDefinition, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, business_id, name, code, description, status
		FROM service_definitions
		WHERE business_id = $1
		ORDER BY name
	`, businessID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	services := []domain.ServiceDefinition{}
	for rows.Next() {
		var service domain.ServiceDefinition
		if err := rows.Scan(&service.ID, &service.BusinessID, &service.Name, &service.Code, &service.Description, &service.Status); err != nil {
			return nil, err
		}
		services = append(services, service)
	}
	return services, rows.Err()
}

func (r *ServiceDefinitionRepository) FindByID(ctx context.Context, serviceDefinitionID string) (domain.ServiceDefinition, error) {
	var service domain.ServiceDefinition
	err := r.db.QueryRow(ctx, `
		SELECT id, business_id, name, code, description, status
		FROM service_definitions
		WHERE id = $1
	`, serviceDefinitionID).Scan(&service.ID, &service.BusinessID, &service.Name, &service.Code, &service.Description, &service.Status)
	return service, err
}
