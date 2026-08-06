package domain

import (
	"time"

	"github.com/google/uuid"
)

type Branch struct {
	ID         uuid.UUID
	BusinessID uuid.UUID
	Name       string
	Code       string
	Address    string
	Phone      string
	Status     string
	CreatedAt  time.Time
}
