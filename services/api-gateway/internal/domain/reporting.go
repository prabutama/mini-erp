package domain

type AuditEventResponse struct {
	EventID      string `json:"event_id"`
	EventType    string `json:"event_type"`
	EventVersion int    `json:"event_version"`
	Producer     string `json:"producer"`
	BusinessID   string `json:"business_id"`
	BranchID     string `json:"branch_id"`
	ActorID      string `json:"actor_id"`
	RequestID    string `json:"request_id"`
	OccurredAt   string `json:"occurred_at"`
	Data         string `json:"data"`
}
type ListAuditEventsResponse struct {
	Events []AuditEventResponse `json:"events"`
}
type OperationsSummaryReportResponse struct {
	BusinessID       string  `json:"business_id"`
	BranchID         string  `json:"branch_id"`
	SnapshotDate     string  `json:"snapshot_date"`
	OpenOrders       int     `json:"open_orders"`
	InProgressOrders int     `json:"in_progress_orders"`
	CompletedOrders  int     `json:"completed_orders"`
	CancelledOrders  int     `json:"cancelled_orders"`
	ResourcesUsed    float64 `json:"resources_used"`
}
