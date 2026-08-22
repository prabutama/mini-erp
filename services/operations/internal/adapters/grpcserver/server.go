package grpcserver

import (
	"context"
	"errors"

	"github.com/isapr/mini-erp/services/operations/internal/application"
	"github.com/isapr/mini-erp/services/operations/internal/domain"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/status"
)

const serviceName = "operations.v1.OperationsService"

type CreateServiceDefinitionRequest struct {
	BusinessID  string `json:"business_id"`
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

type ListServiceDefinitionsRequest struct {
	BusinessID string `json:"business_id"`
}

type ListServiceDefinitionsResponse struct {
	ServiceDefinitions []ServiceDefinitionResponse `json:"service_definitions"`
}

type CreateServiceOrderRequest struct {
	BusinessID          string `json:"business_id"`
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

type ListServiceOrdersRequest struct {
	BusinessID     string `json:"business_id"`
	BranchID       string `json:"branch_id"`
	Status         string `json:"status"`
	AssignedUserID string `json:"assigned_user_id"`
}

type ListServiceOrdersResponse struct {
	ServiceOrders []ServiceOrderResponse `json:"service_orders"`
}

type ServiceOrderSummaryRequest struct {
	BusinessID     string `json:"business_id"`
	BranchID       string `json:"branch_id"`
	AssignedUserID string `json:"assigned_user_id"`
}

type ServiceOrderSummaryResponse struct {
	Total      int `json:"total"`
	Open       int `json:"open"`
	InProgress int `json:"in_progress"`
	Completed  int `json:"completed"`
	Cancelled  int `json:"cancelled"`
}

type GetServiceOrderRequest struct {
	ServiceOrderID string `json:"service_order_id"`
}

type TransitionServiceOrderRequest struct {
	ServiceOrderID  string `json:"service_order_id"`
	BusinessID      string `json:"business_id"`
	Status          string `json:"status"`
	ChangedByUserID string `json:"changed_by_user_id"`
	RequestID       string `json:"request_id"`
}

type AssignServiceOrderRequest struct {
	ServiceOrderID   string `json:"service_order_id"`
	BusinessID       string `json:"business_id"`
	AssignedUserID   string `json:"assigned_user_id"`
	AssignedByUserID string `json:"assigned_by_user_id"`
	RequestID        string `json:"request_id"`
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

type ListAssignedServiceOrdersRequest struct {
	BusinessID     string `json:"business_id"`
	AssignedUserID string `json:"assigned_user_id"`
	BranchID       string `json:"branch_id"`
}

type ListAssignedServiceOrdersResponse struct {
	ServiceOrders []ServiceOrderResponse `json:"service_orders"`
}

type ListServiceOrderAssignmentsRequest struct {
	ServiceOrderID string `json:"service_order_id"`
}

type ListServiceOrderAssignmentsResponse struct {
	Assignments []ServiceOrderAssignmentResponse `json:"assignments"`
}

type CreateWorkflowRequest struct {
	BusinessID  string `json:"business_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
type UpdateWorkflowRequest struct {
	WorkflowID  string `json:"workflow_id"`
	BusinessID  string `json:"business_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
type GetWorkflowRequest struct {
	WorkflowID string `json:"workflow_id"`
}
type ListWorkflowsRequest struct {
	BusinessID string `json:"business_id"`
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
	WorkflowID string `json:"workflow_id"`
	BusinessID string `json:"business_id"`
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
	WorkflowID     string `json:"workflow_id"`
	BusinessID     string `json:"business_id"`
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

type Server struct {
	services  *application.ServiceDefinitionService
	orders    *application.ServiceOrderService
	workflows *application.WorkflowService
}

type operationsServiceServer interface {
	CreateServiceDefinition(context.Context, CreateServiceDefinitionRequest) (ServiceDefinitionResponse, error)
	ListServiceDefinitions(context.Context, ListServiceDefinitionsRequest) (ListServiceDefinitionsResponse, error)
	CreateServiceOrder(context.Context, CreateServiceOrderRequest) (ServiceOrderResponse, error)
	ListServiceOrders(context.Context, ListServiceOrdersRequest) (ListServiceOrdersResponse, error)
	ServiceOrderSummary(context.Context, ServiceOrderSummaryRequest) (ServiceOrderSummaryResponse, error)
	GetServiceOrder(context.Context, GetServiceOrderRequest) (ServiceOrderResponse, error)
	TransitionServiceOrder(context.Context, TransitionServiceOrderRequest) (ServiceOrderResponse, error)
	AssignServiceOrder(context.Context, AssignServiceOrderRequest) (ServiceOrderAssignmentResponse, error)
	ListAssignedServiceOrders(context.Context, ListAssignedServiceOrdersRequest) (ListAssignedServiceOrdersResponse, error)
	ListServiceOrderAssignments(context.Context, ListServiceOrderAssignmentsRequest) (ListServiceOrderAssignmentsResponse, error)
	CreateWorkflow(context.Context, CreateWorkflowRequest) (WorkflowResponse, error)
	ListWorkflows(context.Context, ListWorkflowsRequest) (ListWorkflowsResponse, error)
	GetWorkflow(context.Context, GetWorkflowRequest) (WorkflowResponse, error)
	UpdateWorkflow(context.Context, UpdateWorkflowRequest) (WorkflowResponse, error)
	CreateWorkflowStatus(context.Context, CreateWorkflowStatusRequest) (WorkflowStatusResponse, error)
	CreateWorkflowTransition(context.Context, CreateWorkflowTransitionRequest) (WorkflowTransitionResponse, error)
}

func New(services *application.ServiceDefinitionService, orders *application.ServiceOrderService, workflows *application.WorkflowService) *Server {
	encoding.RegisterCodec(jsonCodec{})
	return &Server{services: services, orders: orders, workflows: workflows}
}

func (s *Server) Register(grpcServer *grpc.Server) {
	grpcServer.RegisterService(&grpc.ServiceDesc{ServiceName: serviceName, HandlerType: (*operationsServiceServer)(nil), Methods: []grpc.MethodDesc{
		{MethodName: "CreateServiceDefinition", Handler: createServiceDefinitionHandler},
		{MethodName: "ListServiceDefinitions", Handler: listServiceDefinitionsHandler},
		{MethodName: "CreateServiceOrder", Handler: createServiceOrderHandler},
		{MethodName: "ListServiceOrders", Handler: listServiceOrdersHandler},
		{MethodName: "ServiceOrderSummary", Handler: serviceOrderSummaryHandler},
		{MethodName: "GetServiceOrder", Handler: getServiceOrderHandler},
		{MethodName: "TransitionServiceOrder", Handler: transitionServiceOrderHandler},
		{MethodName: "AssignServiceOrder", Handler: assignServiceOrderHandler},
		{MethodName: "ListAssignedServiceOrders", Handler: listAssignedServiceOrdersHandler},
		{MethodName: "ListServiceOrderAssignments", Handler: listServiceOrderAssignmentsHandler},
		{MethodName: "CreateWorkflow", Handler: createWorkflowHandler},
		{MethodName: "ListWorkflows", Handler: listWorkflowsHandler},
		{MethodName: "GetWorkflow", Handler: getWorkflowHandler},
		{MethodName: "UpdateWorkflow", Handler: updateWorkflowHandler},
		{MethodName: "CreateWorkflowStatus", Handler: createWorkflowStatusHandler},
		{MethodName: "CreateWorkflowTransition", Handler: createWorkflowTransitionHandler},
	}}, s)
}

func createServiceDefinitionHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req CreateServiceDefinitionRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(operationsServiceServer).CreateServiceDefinition(ctx, req)
}

func listServiceDefinitionsHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req ListServiceDefinitionsRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(operationsServiceServer).ListServiceDefinitions(ctx, req)
}

func createServiceOrderHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req CreateServiceOrderRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(operationsServiceServer).CreateServiceOrder(ctx, req)
}

func listServiceOrdersHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req ListServiceOrdersRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(operationsServiceServer).ListServiceOrders(ctx, req)
}

func serviceOrderSummaryHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req ServiceOrderSummaryRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(operationsServiceServer).ServiceOrderSummary(ctx, req)
}

func getServiceOrderHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req GetServiceOrderRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(operationsServiceServer).GetServiceOrder(ctx, req)
}

func transitionServiceOrderHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req TransitionServiceOrderRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(operationsServiceServer).TransitionServiceOrder(ctx, req)
}

func assignServiceOrderHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req AssignServiceOrderRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(operationsServiceServer).AssignServiceOrder(ctx, req)
}

func listAssignedServiceOrdersHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req ListAssignedServiceOrdersRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(operationsServiceServer).ListAssignedServiceOrders(ctx, req)
}

func listServiceOrderAssignmentsHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req ListServiceOrderAssignmentsRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(operationsServiceServer).ListServiceOrderAssignments(ctx, req)
}

func createWorkflowHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req CreateWorkflowRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(operationsServiceServer).CreateWorkflow(ctx, req)
}
func listWorkflowsHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req ListWorkflowsRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(operationsServiceServer).ListWorkflows(ctx, req)
}
func getWorkflowHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req GetWorkflowRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(operationsServiceServer).GetWorkflow(ctx, req)
}
func updateWorkflowHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req UpdateWorkflowRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(operationsServiceServer).UpdateWorkflow(ctx, req)
}
func createWorkflowStatusHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req CreateWorkflowStatusRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(operationsServiceServer).CreateWorkflowStatus(ctx, req)
}
func createWorkflowTransitionHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req CreateWorkflowTransitionRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(operationsServiceServer).CreateWorkflowTransition(ctx, req)
}

func (s *Server) CreateServiceDefinition(ctx context.Context, req CreateServiceDefinitionRequest) (ServiceDefinitionResponse, error) {
	service, err := s.services.Create(ctx, application.CreateServiceDefinitionInput{BusinessID: req.BusinessID, Name: req.Name, Description: req.Description})
	if err != nil {
		return ServiceDefinitionResponse{}, mapError(err)
	}
	return serviceDefinitionResponse(service), nil
}

func (s *Server) ListServiceDefinitions(ctx context.Context, req ListServiceDefinitionsRequest) (ListServiceDefinitionsResponse, error) {
	services, err := s.services.List(ctx, req.BusinessID)
	if err != nil {
		return ListServiceDefinitionsResponse{}, mapError(err)
	}
	response := ListServiceDefinitionsResponse{ServiceDefinitions: []ServiceDefinitionResponse{}}
	for _, service := range services {
		response.ServiceDefinitions = append(response.ServiceDefinitions, serviceDefinitionResponse(service))
	}
	return response, nil
}

func (s *Server) CreateServiceOrder(ctx context.Context, req CreateServiceOrderRequest) (ServiceOrderResponse, error) {
	order, err := s.orders.Create(ctx, application.CreateServiceOrderInput{BusinessID: req.BusinessID, BranchID: req.BranchID, ServiceDefinitionID: req.ServiceDefinitionID, Title: req.Title, Description: req.Description, Priority: req.Priority})
	if err != nil {
		return ServiceOrderResponse{}, mapError(err)
	}
	return serviceOrderResponse(order), nil
}

func (s *Server) ListServiceOrders(ctx context.Context, req ListServiceOrdersRequest) (ListServiceOrdersResponse, error) {
	orders, err := s.orders.List(ctx, application.ListServiceOrdersInput{BusinessID: req.BusinessID, BranchID: req.BranchID, Status: req.Status, AssignedUserID: req.AssignedUserID})
	if err != nil {
		return ListServiceOrdersResponse{}, mapError(err)
	}
	response := ListServiceOrdersResponse{ServiceOrders: []ServiceOrderResponse{}}
	for _, order := range orders {
		response.ServiceOrders = append(response.ServiceOrders, serviceOrderResponse(order))
	}
	return response, nil
}

func (s *Server) ServiceOrderSummary(ctx context.Context, req ServiceOrderSummaryRequest) (ServiceOrderSummaryResponse, error) {
	summary, err := s.orders.Summary(ctx, application.ServiceOrderSummaryInput{BusinessID: req.BusinessID, BranchID: req.BranchID, AssignedUserID: req.AssignedUserID})
	if err != nil {
		return ServiceOrderSummaryResponse{}, mapError(err)
	}
	return ServiceOrderSummaryResponse{Total: summary.Total, Open: summary.Open, InProgress: summary.InProgress, Completed: summary.Completed, Cancelled: summary.Cancelled}, nil
}

func (s *Server) GetServiceOrder(ctx context.Context, req GetServiceOrderRequest) (ServiceOrderResponse, error) {
	order, err := s.orders.Get(ctx, req.ServiceOrderID)
	if err != nil {
		return ServiceOrderResponse{}, mapError(err)
	}
	return serviceOrderResponse(order), nil
}

func (s *Server) TransitionServiceOrder(ctx context.Context, req TransitionServiceOrderRequest) (ServiceOrderResponse, error) {
	order, err := s.orders.Transition(ctx, application.TransitionServiceOrderInput{ServiceOrderID: req.ServiceOrderID, BusinessID: req.BusinessID, Status: req.Status, ChangedByUserID: req.ChangedByUserID, RequestID: req.RequestID})
	if err != nil {
		return ServiceOrderResponse{}, mapError(err)
	}
	return serviceOrderResponse(order), nil
}

func (s *Server) AssignServiceOrder(ctx context.Context, req AssignServiceOrderRequest) (ServiceOrderAssignmentResponse, error) {
	assignment, err := s.orders.Assign(ctx, application.AssignServiceOrderInput{ServiceOrderID: req.ServiceOrderID, BusinessID: req.BusinessID, AssignedUserID: req.AssignedUserID, AssignedByUserID: req.AssignedByUserID, RequestID: req.RequestID})
	if err != nil {
		return ServiceOrderAssignmentResponse{}, mapError(err)
	}
	return assignmentResponse(assignment), nil
}

func (s *Server) ListAssignedServiceOrders(ctx context.Context, req ListAssignedServiceOrdersRequest) (ListAssignedServiceOrdersResponse, error) {
	orders, err := s.orders.ListAssignedToUser(ctx, req.BusinessID, req.AssignedUserID, req.BranchID)
	if err != nil {
		return ListAssignedServiceOrdersResponse{}, mapError(err)
	}
	response := ListAssignedServiceOrdersResponse{ServiceOrders: []ServiceOrderResponse{}}
	for _, order := range orders {
		response.ServiceOrders = append(response.ServiceOrders, serviceOrderResponse(order))
	}
	return response, nil
}

func (s *Server) ListServiceOrderAssignments(ctx context.Context, req ListServiceOrderAssignmentsRequest) (ListServiceOrderAssignmentsResponse, error) {
	assignments, err := s.orders.ListAssignments(ctx, req.ServiceOrderID)
	if err != nil {
		return ListServiceOrderAssignmentsResponse{}, mapError(err)
	}
	response := ListServiceOrderAssignmentsResponse{Assignments: []ServiceOrderAssignmentResponse{}}
	for _, assignment := range assignments {
		response.Assignments = append(response.Assignments, assignmentResponse(assignment))
	}
	return response, nil
}

func (s *Server) CreateWorkflow(ctx context.Context, req CreateWorkflowRequest) (WorkflowResponse, error) {
	workflow, err := s.workflows.Create(ctx, application.CreateWorkflowInput{BusinessID: req.BusinessID, Name: req.Name, Description: req.Description})
	if err != nil {
		return WorkflowResponse{}, mapError(err)
	}
	return workflowResponse(workflow), nil
}
func (s *Server) ListWorkflows(ctx context.Context, req ListWorkflowsRequest) (ListWorkflowsResponse, error) {
	workflows, err := s.workflows.List(ctx, req.BusinessID)
	if err != nil {
		return ListWorkflowsResponse{}, mapError(err)
	}
	resp := ListWorkflowsResponse{Workflows: []WorkflowResponse{}}
	for _, workflow := range workflows {
		resp.Workflows = append(resp.Workflows, workflowResponse(workflow))
	}
	return resp, nil
}
func (s *Server) GetWorkflow(ctx context.Context, req GetWorkflowRequest) (WorkflowResponse, error) {
	workflow, err := s.workflows.Get(ctx, req.WorkflowID)
	if err != nil {
		return WorkflowResponse{}, mapError(err)
	}
	return workflowResponse(workflow), nil
}
func (s *Server) UpdateWorkflow(ctx context.Context, req UpdateWorkflowRequest) (WorkflowResponse, error) {
	workflow, err := s.workflows.Update(ctx, application.UpdateWorkflowInput{WorkflowID: req.WorkflowID, BusinessID: req.BusinessID, Name: req.Name, Description: req.Description, Status: req.Status})
	if err != nil {
		return WorkflowResponse{}, mapError(err)
	}
	return workflowResponse(workflow), nil
}
func (s *Server) CreateWorkflowStatus(ctx context.Context, req CreateWorkflowStatusRequest) (WorkflowStatusResponse, error) {
	item, err := s.workflows.AddStatus(ctx, application.CreateWorkflowStatusInput{WorkflowID: req.WorkflowID, BusinessID: req.BusinessID, Code: req.Code, Name: req.Name, Category: req.Category, SortOrder: req.SortOrder, IsInitial: req.IsInitial, IsTerminal: req.IsTerminal})
	if err != nil {
		return WorkflowStatusResponse{}, mapError(err)
	}
	return workflowStatusResponse(item), nil
}
func (s *Server) CreateWorkflowTransition(ctx context.Context, req CreateWorkflowTransitionRequest) (WorkflowTransitionResponse, error) {
	item, err := s.workflows.AddTransition(ctx, application.CreateWorkflowTransitionInput{WorkflowID: req.WorkflowID, BusinessID: req.BusinessID, FromStatusCode: req.FromStatusCode, ToStatusCode: req.ToStatusCode})
	if err != nil {
		return WorkflowTransitionResponse{}, mapError(err)
	}
	return workflowTransitionResponse(item), nil
}

func mapError(err error) error {
	if errors.Is(err, application.ErrValidation) {
		return status.Error(codes.InvalidArgument, "validation failed")
	}
	if errors.Is(err, application.ErrInvalidStatusTransition) {
		return status.Error(codes.FailedPrecondition, "invalid status transition")
	}
	if errors.Is(err, application.ErrServiceOrderClosed) {
		return status.Error(codes.FailedPrecondition, "service order closed")
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return status.Error(codes.NotFound, "not found")
	}
	return status.Error(codes.Internal, err.Error())
}

func serviceDefinitionResponse(service domain.ServiceDefinition) ServiceDefinitionResponse {
	return ServiceDefinitionResponse{ServiceDefinitionID: service.ID.String(), BusinessID: service.BusinessID.String(), Name: service.Name, Code: service.Code, Description: service.Description, Status: service.Status}
}

func serviceOrderResponse(order domain.ServiceOrder) ServiceOrderResponse {
	return ServiceOrderResponse{ServiceOrderID: order.ID.String(), BusinessID: order.BusinessID.String(), BranchID: order.BranchID.String(), ServiceDefinitionID: order.ServiceDefinitionID.String(), Title: order.Title, Description: order.Description, Status: order.Status, Priority: order.Priority}
}

func assignmentResponse(assignment domain.ServiceOrderAssignment) ServiceOrderAssignmentResponse {
	return ServiceOrderAssignmentResponse{AssignmentID: assignment.ID.String(), ServiceOrderID: assignment.ServiceOrderID.String(), BusinessID: assignment.BusinessID.String(), BranchID: assignment.BranchID.String(), AssignedUserID: assignment.AssignedUserID.String(), AssignedByUserID: assignment.AssignedByUserID.String(), Status: assignment.Status}
}

func workflowResponse(workflow domain.Workflow) WorkflowResponse {
	return WorkflowResponse{WorkflowID: workflow.ID.String(), BusinessID: workflow.BusinessID.String(), Name: workflow.Name, Description: workflow.Description, Status: workflow.Status}
}
func workflowStatusResponse(status domain.WorkflowStatus) WorkflowStatusResponse {
	return WorkflowStatusResponse{WorkflowStatusID: status.ID.String(), WorkflowID: status.WorkflowID.String(), BusinessID: status.BusinessID.String(), Code: status.Code, Name: status.Name, Category: status.Category, SortOrder: status.SortOrder, IsInitial: status.IsInitial, IsTerminal: status.IsTerminal}
}
func workflowTransitionResponse(transition domain.WorkflowTransition) WorkflowTransitionResponse {
	return WorkflowTransitionResponse{WorkflowTransitionID: transition.ID.String(), WorkflowID: transition.WorkflowID.String(), BusinessID: transition.BusinessID.String(), FromStatusCode: transition.FromStatusCode, ToStatusCode: transition.ToStatusCode}
}
