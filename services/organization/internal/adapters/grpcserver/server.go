package grpcserver

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/isapr/mini-erp/services/organization/internal/application"
	"github.com/isapr/mini-erp/services/organization/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/status"
)

const serviceName = "organization.v1.OrganizationService"

type CreateBusinessRequest struct {
	Name     string `json:"name"`
	Timezone string `json:"timezone"`
}

type CreateBusinessResponse struct {
	BusinessID    string `json:"business_id"`
	Name          string `json:"name"`
	Code          string `json:"code"`
	Status        string `json:"status"`
	Plan          string `json:"plan"`
	PlatformNotes string `json:"platform_notes"`
	SuspendedAt   string `json:"suspended_at"`
	Timezone      string `json:"timezone"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type GetBusinessRequest struct {
	BusinessID string `json:"business_id"`
}

type ListPlatformBusinessesRequest struct{}

type ListPlatformBusinessesResponse struct {
	Businesses []CreateBusinessResponse `json:"businesses"`
}

type UpdateBusinessRequest struct {
	BusinessID string `json:"business_id"`
	Name       string `json:"name"`
	Timezone   string `json:"timezone"`
}

type UpdatePlatformBusinessRequest struct {
	BusinessID    string `json:"business_id"`
	Status        string `json:"status"`
	Plan          string `json:"plan"`
	PlatformNotes string `json:"platform_notes"`
	SuspendedAt   string `json:"suspended_at"`
}

type CreateBranchRequest struct {
	BusinessID string `json:"business_id"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	Phone      string `json:"phone"`
}

type BranchResponse struct {
	BranchID   string `json:"branch_id"`
	BusinessID string `json:"business_id"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	Address    string `json:"address"`
	Phone      string `json:"phone"`
	Status     string `json:"status"`
}

type GetBranchRequest struct {
	BranchID string `json:"branch_id"`
}

type ListBranchesRequest struct {
	BusinessID string `json:"business_id"`
}

type ListBranchesResponse struct {
	Branches []BranchResponse `json:"branches"`
}

type UpdateBranchRequest struct {
	BranchID string `json:"branch_id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Status   string `json:"status"`
}

type CreateEmployeePlacementRequest struct {
	UserID         string `json:"user_id"`
	BusinessID     string `json:"business_id"`
	BranchID       string `json:"branch_id"`
	Position       string `json:"position"`
	EmploymentType string `json:"employment_type"`
}

type EmployeePlacementResponse struct {
	PlacementID    string `json:"placement_id"`
	UserID         string `json:"user_id"`
	BusinessID     string `json:"business_id"`
	BranchID       string `json:"branch_id"`
	Position       string `json:"position"`
	EmploymentType string `json:"employment_type"`
	Status         string `json:"status"`
}

type ListAssignedBranchesRequest struct {
	UserID     string `json:"user_id"`
	BusinessID string `json:"business_id"`
}

type ListAssignedBranchesResponse struct {
	BranchIDs  []string                    `json:"branch_ids"`
	Placements []EmployeePlacementResponse `json:"placements"`
}

type Server struct {
	businesses *application.BusinessService
	branches   *application.BranchService
	placements *application.PlacementService
}

type organizationServiceServer interface {
	CreateBusiness(context.Context, CreateBusinessRequest) (CreateBusinessResponse, error)
	GetBusiness(context.Context, GetBusinessRequest) (CreateBusinessResponse, error)
	UpdateBusiness(context.Context, UpdateBusinessRequest) (CreateBusinessResponse, error)
	ListPlatformBusinesses(context.Context, ListPlatformBusinessesRequest) (ListPlatformBusinessesResponse, error)
	UpdatePlatformBusiness(context.Context, UpdatePlatformBusinessRequest) (CreateBusinessResponse, error)
	CreateBranch(context.Context, CreateBranchRequest) (BranchResponse, error)
	GetBranch(context.Context, GetBranchRequest) (BranchResponse, error)
	ListBranches(context.Context, ListBranchesRequest) (ListBranchesResponse, error)
	UpdateBranch(context.Context, UpdateBranchRequest) (BranchResponse, error)
	CreateEmployeePlacement(context.Context, CreateEmployeePlacementRequest) (EmployeePlacementResponse, error)
	ListAssignedBranches(context.Context, ListAssignedBranchesRequest) (ListAssignedBranchesResponse, error)
}

func New(businesses *application.BusinessService, branches *application.BranchService, placements *application.PlacementService) *Server {
	encoding.RegisterCodec(jsonCodec{})
	return &Server{businesses: businesses, branches: branches, placements: placements}
}

func (s *Server) Register(grpcServer *grpc.Server) {
	grpcServer.RegisterService(&grpc.ServiceDesc{
		ServiceName: serviceName,
		HandlerType: (*organizationServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			{MethodName: "CreateBusiness", Handler: createBusinessHandler},
			{MethodName: "GetBusiness", Handler: getBusinessHandler},
			{MethodName: "UpdateBusiness", Handler: updateBusinessHandler},
			{MethodName: "ListPlatformBusinesses", Handler: listPlatformBusinessesHandler},
			{MethodName: "UpdatePlatformBusiness", Handler: updatePlatformBusinessHandler},
			{MethodName: "CreateBranch", Handler: createBranchHandler},
			{MethodName: "GetBranch", Handler: getBranchHandler},
			{MethodName: "ListBranches", Handler: listBranchesHandler},
			{MethodName: "UpdateBranch", Handler: updateBranchHandler},
			{MethodName: "CreateEmployeePlacement", Handler: createEmployeePlacementHandler},
			{MethodName: "ListAssignedBranches", Handler: listAssignedBranchesHandler},
		},
	}, s)
}

func createBusinessHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req CreateBusinessRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(organizationServiceServer).CreateBusiness(ctx, req)
}

func getBusinessHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req GetBusinessRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(organizationServiceServer).GetBusiness(ctx, req)
}

func updateBusinessHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req UpdateBusinessRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(organizationServiceServer).UpdateBusiness(ctx, req)
}

func listPlatformBusinessesHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req ListPlatformBusinessesRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(organizationServiceServer).ListPlatformBusinesses(ctx, req)
}

func updatePlatformBusinessHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req UpdatePlatformBusinessRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(organizationServiceServer).UpdatePlatformBusiness(ctx, req)
}

func createBranchHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req CreateBranchRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(organizationServiceServer).CreateBranch(ctx, req)
}

func getBranchHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req GetBranchRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(organizationServiceServer).GetBranch(ctx, req)
}

func listBranchesHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req ListBranchesRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(organizationServiceServer).ListBranches(ctx, req)
}

func updateBranchHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req UpdateBranchRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(organizationServiceServer).UpdateBranch(ctx, req)
}

func createEmployeePlacementHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req CreateEmployeePlacementRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(organizationServiceServer).CreateEmployeePlacement(ctx, req)
}

func listAssignedBranchesHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req ListAssignedBranchesRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(organizationServiceServer).ListAssignedBranches(ctx, req)
}

func (s *Server) CreateBusiness(ctx context.Context, req CreateBusinessRequest) (CreateBusinessResponse, error) {
	business, err := s.businesses.CreateBusiness(ctx, application.CreateBusinessInput{Name: req.Name, Timezone: req.Timezone})
	if errors.Is(err, application.ErrValidation) {
		return CreateBusinessResponse{}, status.Error(codes.InvalidArgument, "validation failed")
	}
	if errors.Is(err, application.ErrBusinessAlreadyExists) {
		return CreateBusinessResponse{}, status.Error(codes.AlreadyExists, "business already exists")
	}
	if err != nil {
		return CreateBusinessResponse{}, status.Error(codes.Internal, err.Error())
	}

	log.Printf("grpc service=organization method=CreateBusiness business_id=%s code=%s status=%s", business.ID, business.Code, business.Status)

	return businessResponse(business), nil
}

func (s *Server) GetBusiness(ctx context.Context, req GetBusinessRequest) (CreateBusinessResponse, error) {
	business, err := s.businesses.GetBusiness(ctx, req.BusinessID)
	if err != nil {
		return CreateBusinessResponse{}, mapBusinessError(err)
	}
	return businessResponse(business), nil
}

func (s *Server) UpdateBusiness(ctx context.Context, req UpdateBusinessRequest) (CreateBusinessResponse, error) {
	business, err := s.businesses.UpdateBusiness(ctx, application.UpdateBusinessInput{BusinessID: req.BusinessID, Name: req.Name, Timezone: req.Timezone})
	if err != nil {
		return CreateBusinessResponse{}, mapBusinessError(err)
	}
	return businessResponse(business), nil
}

func (s *Server) ListPlatformBusinesses(ctx context.Context, _ ListPlatformBusinessesRequest) (ListPlatformBusinessesResponse, error) {
	businesses, err := s.businesses.ListPlatformBusinesses(ctx)
	if err != nil {
		return ListPlatformBusinessesResponse{}, status.Error(codes.Internal, err.Error())
	}
	response := ListPlatformBusinessesResponse{Businesses: []CreateBusinessResponse{}}
	for _, business := range businesses {
		response.Businesses = append(response.Businesses, businessResponse(business))
	}
	return response, nil
}

func (s *Server) UpdatePlatformBusiness(ctx context.Context, req UpdatePlatformBusinessRequest) (CreateBusinessResponse, error) {
	business, err := s.businesses.UpdatePlatformBusiness(ctx, application.UpdatePlatformBusinessInput{BusinessID: req.BusinessID, Status: req.Status, Plan: req.Plan, PlatformNotes: req.PlatformNotes, SuspendedAtRFC: req.SuspendedAt})
	if err != nil {
		return CreateBusinessResponse{}, mapBusinessError(err)
	}
	return businessResponse(business), nil
}

func (s *Server) CreateBranch(ctx context.Context, req CreateBranchRequest) (BranchResponse, error) {
	branch, err := s.branches.CreateBranch(ctx, application.CreateBranchInput{BusinessID: req.BusinessID, Name: req.Name, Address: req.Address, Phone: req.Phone})
	if errors.Is(err, application.ErrValidation) {
		return BranchResponse{}, status.Error(codes.InvalidArgument, "validation failed")
	}
	if err != nil {
		return BranchResponse{}, status.Error(codes.Internal, err.Error())
	}
	log.Printf("grpc service=organization method=CreateBranch branch_id=%s business_id=%s status=%s", branch.ID, branch.BusinessID, branch.Status)
	return branchResponse(branch), nil
}

func (s *Server) GetBranch(ctx context.Context, req GetBranchRequest) (BranchResponse, error) {
	branch, err := s.branches.GetBranch(ctx, req.BranchID)
	if err != nil {
		return BranchResponse{}, status.Error(codes.NotFound, "branch not found")
	}
	return branchResponse(branch), nil
}

func (s *Server) ListBranches(ctx context.Context, req ListBranchesRequest) (ListBranchesResponse, error) {
	branches, err := s.branches.ListBranches(ctx, req.BusinessID)
	if errors.Is(err, application.ErrValidation) {
		return ListBranchesResponse{}, status.Error(codes.InvalidArgument, "validation failed")
	}
	if err != nil {
		return ListBranchesResponse{}, status.Error(codes.Internal, err.Error())
	}
	response := ListBranchesResponse{Branches: []BranchResponse{}}
	for _, branch := range branches {
		response.Branches = append(response.Branches, branchResponse(branch))
	}
	return response, nil
}

func (s *Server) UpdateBranch(ctx context.Context, req UpdateBranchRequest) (BranchResponse, error) {
	branch, err := s.branches.UpdateBranch(ctx, application.UpdateBranchInput{BranchID: req.BranchID, Name: req.Name, Address: req.Address, Phone: req.Phone, Status: req.Status})
	if err != nil {
		return BranchResponse{}, status.Error(codes.Internal, err.Error())
	}
	return branchResponse(branch), nil
}

func (s *Server) CreateEmployeePlacement(ctx context.Context, req CreateEmployeePlacementRequest) (EmployeePlacementResponse, error) {
	placement, err := s.placements.CreatePlacement(ctx, application.CreatePlacementInput{UserID: req.UserID, BusinessID: req.BusinessID, BranchID: req.BranchID, Position: req.Position, EmploymentType: req.EmploymentType})
	if errors.Is(err, application.ErrValidation) {
		return EmployeePlacementResponse{}, status.Error(codes.InvalidArgument, "validation failed")
	}
	if err != nil {
		return EmployeePlacementResponse{}, status.Error(codes.Internal, err.Error())
	}
	return EmployeePlacementResponse{PlacementID: placement.ID.String(), UserID: placement.UserID.String(), BusinessID: placement.BusinessID.String(), BranchID: placement.BranchID.String(), Position: placement.Position, EmploymentType: placement.EmploymentType, Status: placement.Status}, nil
}

func (s *Server) ListAssignedBranches(ctx context.Context, req ListAssignedBranchesRequest) (ListAssignedBranchesResponse, error) {
	placements, err := s.placements.ListAssignedBranches(ctx, application.ListAssignedBranchesInput{UserID: req.UserID, BusinessID: req.BusinessID})
	if errors.Is(err, application.ErrValidation) {
		return ListAssignedBranchesResponse{}, status.Error(codes.InvalidArgument, "validation failed")
	}
	if err != nil {
		return ListAssignedBranchesResponse{}, status.Error(codes.Internal, err.Error())
	}

	response := ListAssignedBranchesResponse{BranchIDs: []string{}, Placements: []EmployeePlacementResponse{}}
	for _, placement := range placements {
		response.BranchIDs = append(response.BranchIDs, placement.BranchID.String())
		response.Placements = append(response.Placements, EmployeePlacementResponse{PlacementID: placement.ID.String(), UserID: placement.UserID.String(), BusinessID: placement.BusinessID.String(), BranchID: placement.BranchID.String(), Position: placement.Position, EmploymentType: placement.EmploymentType, Status: placement.Status})
	}
	return response, nil
}

func branchResponse(branch domain.Branch) BranchResponse {
	return BranchResponse{BranchID: branch.ID.String(), BusinessID: branch.BusinessID.String(), Name: branch.Name, Code: branch.Code, Address: branch.Address, Phone: branch.Phone, Status: branch.Status}
}

func businessResponse(business domain.Business) CreateBusinessResponse {
	suspendedAt := ""
	if business.SuspendedAt != nil {
		suspendedAt = business.SuspendedAt.Format(time.RFC3339)
	}
	return CreateBusinessResponse{BusinessID: business.ID.String(), Name: business.Name, Code: business.Code, Status: business.Status, Plan: business.Plan, PlatformNotes: business.PlatformNotes, SuspendedAt: suspendedAt, Timezone: business.Timezone, CreatedAt: business.CreatedAt.Format(time.RFC3339), UpdatedAt: business.UpdatedAt.Format(time.RFC3339)}
}

func mapBusinessError(err error) error {
	if errors.Is(err, application.ErrValidation) {
		return status.Error(codes.InvalidArgument, "validation failed")
	}
	return status.Error(codes.Internal, err.Error())
}
