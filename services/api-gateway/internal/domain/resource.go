package domain

type CreateResourceRequest struct {
	BranchID string `json:"branch_id"`
	Name     string `json:"name"`
	Unit     string `json:"unit"`
	Type     string `json:"type"`
}

type ResourceResponse struct {
	ResourceID string `json:"resource_id"`
	BusinessID string `json:"business_id"`
	BranchID   string `json:"branch_id"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	Unit       string `json:"unit"`
	Type       string `json:"type"`
	Status     string `json:"status"`
}

type ListResourcesResponse struct {
	Resources []ResourceResponse `json:"resources"`
}

type RecordStockMovementRequest struct {
	MovementType   string  `json:"movement_type"`
	Quantity       float64 `json:"quantity"`
	Reason         string  `json:"reason"`
	ServiceOrderID string  `json:"service_order_id"`
}

type StockMovementResponse struct {
	StockMovementID string  `json:"stock_movement_id"`
	BusinessID      string  `json:"business_id"`
	BranchID        string  `json:"branch_id"`
	ResourceID      string  `json:"resource_id"`
	MovementType    string  `json:"movement_type"`
	Reason          string  `json:"reason"`
	ServiceOrderID  string  `json:"service_order_id"`
	ActorUserID     string  `json:"actor_user_id"`
	Quantity        float64 `json:"quantity"`
}

type ResourceAvailabilityResponse struct {
	ResourceID string  `json:"resource_id"`
	Quantity   float64 `json:"quantity"`
}

type RecordResourceUsageRequest struct {
	ResourceID string  `json:"resource_id"`
	Quantity   float64 `json:"quantity"`
	Reason     string  `json:"reason"`
}

type ResourceUsageResponse struct {
	ResourceUsageID  string  `json:"resource_usage_id"`
	BusinessID       string  `json:"business_id"`
	BranchID         string  `json:"branch_id"`
	ServiceOrderID   string  `json:"service_order_id"`
	ResourceID       string  `json:"resource_id"`
	Quantity         float64 `json:"quantity"`
	Reason           string  `json:"reason"`
	RecordedByUserID string  `json:"recorded_by_user_id"`
	StockMovementID  string  `json:"stock_movement_id"`
}

type ListResourceUsageResponse struct {
	Usages []ResourceUsageResponse `json:"usages"`
}
