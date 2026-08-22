package grpcserver

import (
	"context"
	"errors"
	"time"

	"github.com/isapr/mini-erp/services/reporting/internal/application"
	"github.com/isapr/mini-erp/services/reporting/internal/domain"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/status"
)

const serviceName = "reporting.v1.ReportingService"

type RecordAuditEventRequest struct {
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
type GetAuditEventsRequest struct {
	BusinessID string `json:"business_id"`
	BranchID   string `json:"branch_id"`
}
type GetAuditEventsResponse struct {
	Events []AuditEventResponse `json:"events"`
}
type UpsertOperationsSummaryRequest struct {
	BusinessID       string  `json:"business_id"`
	BranchID         string  `json:"branch_id"`
	SnapshotDate     string  `json:"snapshot_date"`
	OpenOrders       int     `json:"open_orders"`
	InProgressOrders int     `json:"in_progress_orders"`
	CompletedOrders  int     `json:"completed_orders"`
	CancelledOrders  int     `json:"cancelled_orders"`
	ResourcesUsed    float64 `json:"resources_used"`
}
type GetOperationsSummaryRequest struct {
	BusinessID   string `json:"business_id"`
	BranchID     string `json:"branch_id"`
	SnapshotDate string `json:"snapshot_date"`
}
type OperationsSummaryResponse struct {
	BusinessID       string  `json:"business_id"`
	BranchID         string  `json:"branch_id"`
	SnapshotDate     string  `json:"snapshot_date"`
	OpenOrders       int     `json:"open_orders"`
	InProgressOrders int     `json:"in_progress_orders"`
	CompletedOrders  int     `json:"completed_orders"`
	CancelledOrders  int     `json:"cancelled_orders"`
	ResourcesUsed    float64 `json:"resources_used"`
}

type Server struct{ reports *application.ReportingService }

type reportingServiceServer interface {
	RecordAuditEvent(context.Context, RecordAuditEventRequest) (AuditEventResponse, error)
	GetAuditEvents(context.Context, GetAuditEventsRequest) (GetAuditEventsResponse, error)
	UpsertOperationsSummary(context.Context, UpsertOperationsSummaryRequest) (OperationsSummaryResponse, error)
	GetOperationsSummary(context.Context, GetOperationsSummaryRequest) (OperationsSummaryResponse, error)
}

func New(reports *application.ReportingService) *Server {
	encoding.RegisterCodec(jsonCodec{})
	return &Server{reports: reports}
}
func (s *Server) Register(grpcServer *grpc.Server) {
	grpcServer.RegisterService(&grpc.ServiceDesc{ServiceName: serviceName, HandlerType: (*reportingServiceServer)(nil), Methods: []grpc.MethodDesc{{MethodName: "RecordAuditEvent", Handler: recordAuditEventHandler}, {MethodName: "GetAuditEvents", Handler: getAuditEventsHandler}, {MethodName: "UpsertOperationsSummary", Handler: upsertOperationsSummaryHandler}, {MethodName: "GetOperationsSummary", Handler: getOperationsSummaryHandler}}}, s)
}
func recordAuditEventHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req RecordAuditEventRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(reportingServiceServer).RecordAuditEvent(ctx, req)
}
func getAuditEventsHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req GetAuditEventsRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(reportingServiceServer).GetAuditEvents(ctx, req)
}
func upsertOperationsSummaryHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req UpsertOperationsSummaryRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(reportingServiceServer).UpsertOperationsSummary(ctx, req)
}
func getOperationsSummaryHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req GetOperationsSummaryRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(reportingServiceServer).GetOperationsSummary(ctx, req)
}

func (s *Server) RecordAuditEvent(ctx context.Context, req RecordAuditEventRequest) (AuditEventResponse, error) {
	occurredAt := time.Time{}
	if req.OccurredAt != "" {
		parsed, err := time.Parse(time.RFC3339, req.OccurredAt)
		if err != nil {
			return AuditEventResponse{}, status.Error(codes.InvalidArgument, "validation failed")
		}
		occurredAt = parsed
	}
	event, err := s.reports.RecordAuditEvent(ctx, application.RecordAuditEventInput{EventType: req.EventType, EventVersion: req.EventVersion, Producer: req.Producer, BusinessID: req.BusinessID, BranchID: req.BranchID, ActorID: req.ActorID, RequestID: req.RequestID, OccurredAt: occurredAt, Data: req.Data})
	if err != nil {
		return AuditEventResponse{}, mapError(err)
	}
	return auditEventResponse(event), nil
}
func (s *Server) GetAuditEvents(ctx context.Context, req GetAuditEventsRequest) (GetAuditEventsResponse, error) {
	events, err := s.reports.ListAuditEvents(ctx, req.BusinessID, req.BranchID)
	if err != nil {
		return GetAuditEventsResponse{}, mapError(err)
	}
	resp := GetAuditEventsResponse{Events: []AuditEventResponse{}}
	for _, event := range events {
		resp.Events = append(resp.Events, auditEventResponse(event))
	}
	return resp, nil
}
func (s *Server) UpsertOperationsSummary(ctx context.Context, req UpsertOperationsSummaryRequest) (OperationsSummaryResponse, error) {
	summary, err := s.reports.UpsertOperationsSummary(ctx, application.OperationsSummaryInput{BusinessID: req.BusinessID, BranchID: req.BranchID, SnapshotDate: req.SnapshotDate, OpenOrders: req.OpenOrders, InProgressOrders: req.InProgressOrders, CompletedOrders: req.CompletedOrders, CancelledOrders: req.CancelledOrders, ResourcesUsed: req.ResourcesUsed})
	if err != nil {
		return OperationsSummaryResponse{}, mapError(err)
	}
	return summaryResponse(summary), nil
}
func (s *Server) GetOperationsSummary(ctx context.Context, req GetOperationsSummaryRequest) (OperationsSummaryResponse, error) {
	summary, err := s.reports.GetOperationsSummary(ctx, req.BusinessID, req.BranchID, req.SnapshotDate)
	if err != nil {
		return OperationsSummaryResponse{}, mapError(err)
	}
	return summaryResponse(summary), nil
}

func mapError(err error) error {
	if errors.Is(err, application.ErrValidation) {
		return status.Error(codes.InvalidArgument, "validation failed")
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return status.Error(codes.NotFound, "not found")
	}
	return status.Error(codes.Internal, err.Error())
}
func auditEventResponse(event domain.AuditEvent) AuditEventResponse {
	return AuditEventResponse{EventID: event.ID.String(), EventType: event.EventType, EventVersion: event.EventVersion, Producer: event.Producer, BusinessID: event.BusinessID.String(), BranchID: event.BranchID.String(), ActorID: event.ActorID.String(), RequestID: event.RequestID, OccurredAt: event.OccurredAt.Format(time.RFC3339), Data: event.Data}
}
func summaryResponse(summary domain.OperationsSummary) OperationsSummaryResponse {
	return OperationsSummaryResponse{BusinessID: summary.BusinessID.String(), BranchID: summary.BranchID.String(), SnapshotDate: summary.SnapshotDate.Format("2006-01-02"), OpenOrders: summary.OpenOrders, InProgressOrders: summary.InProgressOrders, CompletedOrders: summary.CompletedOrders, CancelledOrders: summary.CancelledOrders, ResourcesUsed: summary.ResourcesUsed}
}
