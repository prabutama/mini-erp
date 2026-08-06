package domain

import (
	"time"

	"github.com/google/uuid"
)

type EmployeePlacement struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	BusinessID     uuid.UUID
	BranchID       uuid.UUID
	Position       string
	EmploymentType string
	Status         string
	StartDate      time.Time
}
