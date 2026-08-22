package domain

import "github.com/google/uuid"

type ServiceOrder struct {
	ID                  uuid.UUID
	BusinessID          uuid.UUID
	BranchID            uuid.UUID
	ServiceDefinitionID uuid.UUID
	Title               string
	Description         string
	Status              string
	Priority            string
}

type ServiceOrderAssignment struct {
	ID               uuid.UUID
	ServiceOrderID   uuid.UUID
	BusinessID       uuid.UUID
	BranchID         uuid.UUID
	AssignedUserID   uuid.UUID
	AssignedByUserID uuid.UUID
	Status           string
	RequestID        string
}

type ServiceOrderSummary struct {
	Total      int
	Open       int
	InProgress int
	Completed  int
	Cancelled  int
}
