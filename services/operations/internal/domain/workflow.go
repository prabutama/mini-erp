package domain

import "github.com/google/uuid"

type Workflow struct {
	ID          uuid.UUID
	BusinessID  uuid.UUID
	Name        string
	Description string
	Status      string
}

type WorkflowStatus struct {
	ID         uuid.UUID
	WorkflowID uuid.UUID
	BusinessID uuid.UUID
	Code       string
	Name       string
	Category   string
	SortOrder  int
	IsInitial  bool
	IsTerminal bool
}

type WorkflowTransition struct {
	ID             uuid.UUID
	WorkflowID     uuid.UUID
	BusinessID     uuid.UUID
	FromStatusCode string
	ToStatusCode   string
}
