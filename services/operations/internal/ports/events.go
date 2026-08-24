package ports

import (
	"context"

	"github.com/isapr/mini-erp/services/operations/internal/domain"
)

type EventPublisher interface {
	Publish(ctx context.Context, event domain.EventEnvelope) error
}
