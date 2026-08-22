package domain

import "github.com/google/uuid"

type ServiceDefinition struct {
	ID          uuid.UUID
	BusinessID  uuid.UUID
	Name        string
	Code        string
	Description string
	Status      string
}
