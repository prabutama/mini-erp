package domain

import (
	"time"

	"github.com/google/uuid"
)

type AuditEvent struct {
	ID           uuid.UUID
	EventType    string
	EventVersion int
	Producer     string
	BusinessID   uuid.UUID
	BranchID     uuid.UUID
	ActorID      uuid.UUID
	RequestID    string
	OccurredAt   time.Time
	Data         string
}

type OperationsSummary struct {
	ID               uuid.UUID
	BusinessID       uuid.UUID
	BranchID         uuid.UUID
	SnapshotDate     time.Time
	OpenOrders       int
	InProgressOrders int
	CompletedOrders  int
	CancelledOrders  int
	ResourcesUsed    float64
}
