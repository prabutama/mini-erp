package domain

import "time"

type EventEnvelope struct {
	EventID      string         `json:"event_id"`
	EventType    string         `json:"event_type"`
	EventVersion int            `json:"event_version"`
	OccurredAt   time.Time      `json:"occurred_at"`
	Producer     string         `json:"producer"`
	BusinessID   string         `json:"business_id"`
	BranchID     string         `json:"branch_id,omitempty"`
	ActorID      string         `json:"actor_id,omitempty"`
	RequestID    string         `json:"request_id,omitempty"`
	Data         map[string]any `json:"data"`
}
