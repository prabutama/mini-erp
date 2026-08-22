package domain

import "github.com/google/uuid"

type Resource struct {
	ID         uuid.UUID
	BusinessID uuid.UUID
	BranchID   uuid.UUID
	Name       string
	Code       string
	Unit       string
	Type       string
	Status     string
}

type StockMovement struct {
	ID             uuid.UUID
	BusinessID     uuid.UUID
	BranchID       uuid.UUID
	ResourceID     uuid.UUID
	MovementType   string
	Quantity       float64
	Reason         string
	ServiceOrderID uuid.UUID
	ActorUserID    uuid.UUID
	RequestID      string
}

type ResourceAvailability struct {
	ResourceID uuid.UUID
	Quantity   float64
}

type ResourceUsage struct {
	ID               uuid.UUID
	BusinessID       uuid.UUID
	BranchID         uuid.UUID
	ServiceOrderID   uuid.UUID
	ResourceID       uuid.UUID
	Quantity         float64
	Reason           string
	RecordedByUserID uuid.UUID
	StockMovementID  uuid.UUID
	RequestID        string
}
