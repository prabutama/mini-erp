package domain

type CreateServiceDefinitionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ServiceDefinitionResponse struct {
	ServiceDefinitionID string `json:"service_definition_id"`
	BusinessID          string `json:"business_id"`
	Name                string `json:"name"`
	Code                string `json:"code"`
	Description         string `json:"description"`
	Status              string `json:"status"`
}

type ListServiceDefinitionsResponse struct {
	ServiceDefinitions []ServiceDefinitionResponse `json:"service_definitions"`
}

type CreateServiceOrderRequest struct {
	BranchID            string `json:"branch_id"`
	ServiceDefinitionID string `json:"service_definition_id"`
	Title               string `json:"title"`
	Description         string `json:"description"`
	Priority            string `json:"priority"`
}

type ServiceOrderResponse struct {
	ServiceOrderID      string `json:"service_order_id"`
	BusinessID          string `json:"business_id"`
	BranchID            string `json:"branch_id"`
	ServiceDefinitionID string `json:"service_definition_id"`
	Title               string `json:"title"`
	Description         string `json:"description"`
	Status              string `json:"status"`
	Priority            string `json:"priority"`
}

type ListServiceOrdersResponse struct {
	ServiceOrders []ServiceOrderResponse `json:"service_orders"`
}

type ServiceOrderSummaryResponse struct {
	Total      int `json:"total"`
	Open       int `json:"open"`
	InProgress int `json:"in_progress"`
	Completed  int `json:"completed"`
	Cancelled  int `json:"cancelled"`
}

type TransitionServiceOrderRequest struct {
	Status string `json:"status"`
}

type AssignServiceOrderRequest struct {
	AssignedUserID string `json:"assigned_user_id"`
}

type ServiceOrderAssignmentResponse struct {
	AssignmentID     string `json:"assignment_id"`
	ServiceOrderID   string `json:"service_order_id"`
	BusinessID       string `json:"business_id"`
	BranchID         string `json:"branch_id"`
	AssignedUserID   string `json:"assigned_user_id"`
	AssignedByUserID string `json:"assigned_by_user_id"`
	Status           string `json:"status"`
}

type ListServiceOrderAssignmentsResponse struct {
	Assignments []ServiceOrderAssignmentResponse `json:"assignments"`
}

type CreateWorkflowRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
type UpdateWorkflowRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
type WorkflowResponse struct {
	WorkflowID  string `json:"workflow_id"`
	BusinessID  string `json:"business_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
type ListWorkflowsResponse struct {
	Workflows []WorkflowResponse `json:"workflows"`
}
type CreateWorkflowStatusRequest struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	Category   string `json:"category"`
	SortOrder  int    `json:"sort_order"`
	IsInitial  bool   `json:"is_initial"`
	IsTerminal bool   `json:"is_terminal"`
}
type WorkflowStatusResponse struct {
	WorkflowStatusID string `json:"workflow_status_id"`
	WorkflowID       string `json:"workflow_id"`
	BusinessID       string `json:"business_id"`
	Code             string `json:"code"`
	Name             string `json:"name"`
	Category         string `json:"category"`
	SortOrder        int    `json:"sort_order"`
	IsInitial        bool   `json:"is_initial"`
	IsTerminal       bool   `json:"is_terminal"`
}
type CreateWorkflowTransitionRequest struct {
	FromStatusCode string `json:"from_status_code"`
	ToStatusCode   string `json:"to_status_code"`
}
type WorkflowTransitionResponse struct {
	WorkflowTransitionID string `json:"workflow_transition_id"`
	WorkflowID           string `json:"workflow_id"`
	BusinessID           string `json:"business_id"`
	FromStatusCode       string `json:"from_status_code"`
	ToStatusCode         string `json:"to_status_code"`
}
