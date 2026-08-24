package events

import (
	"context"

	"github.com/isapr/mini-erp/services/operations/internal/domain"
)

type NoopPublisher struct{}

func NewNoopPublisher() *NoopPublisher { return &NoopPublisher{} }

func (p *NoopPublisher) Publish(ctx context.Context, event domain.EventEnvelope) error { return nil }
