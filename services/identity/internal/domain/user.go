package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	FullName     string
	Status       string
	CreatedAt    time.Time
}

type AuthContext struct {
	UserID            uuid.UUID
	BusinessID        uuid.UUID
	Role              RoleName
	Permissions       []string
	AssignedBranchIDs []string
	RequestID         string
}

type AuthSession struct {
	AccessToken  string
	RefreshToken string
	Context      AuthContext
}
