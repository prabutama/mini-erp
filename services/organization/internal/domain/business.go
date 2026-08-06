package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrBusinessAlreadyExists = errors.New("business already exists")

type Business struct {
	ID        uuid.UUID
	Name      string
	Code      string
	Status    string
	Plan      string
	Timezone  string
	CreatedAt time.Time
}
