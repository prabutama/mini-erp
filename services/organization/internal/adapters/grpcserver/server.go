package grpcserver

import (
	"context"
	"errors"
	"log"

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
	BusinessID string `json:"business_id"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	Status     string `json:"status"`
	Plan       string `json:"plan"`
	Timezone   string `json:"timezone"`
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

type Server struct {
	businesses *application.BusinessService
	branches   *application.BranchService
	placements *application.PlacementService
}

type organizationServiceServer interface {
	CreateBusiness(context.Context, CreateBusinessRequest) (CreateBusinessResponse, error)
	CreateBranch(context.Context, CreateBranchRequest) (BranchResponse, error)
	GetBranch(context.Context, GetBranchRequest) (BranchResponse, error)
	ListBranches(context.Context, ListBranchesRequest) (ListBranchesResponse, error)
	UpdateBranch(context.Context, UpdateBranchRequest) (BranchResponse, error)
	CreateEmployeePlacement(context.Context, CreateEmployeePlacementRequest) (EmployeePlacementResponse, error)
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
			{MethodName: "CreateBranch", Handler: createBranchHandler},
			{MethodName: "GetBranch", Handler: getBranchHandler},
			{MethodName: "ListBranches", Handler: listBranchesHandler},
			{MethodName: "UpdateBranch", Handler: updateBranchHandler},
			{MethodName: "CreateEmployeePlacement", Handler: createEmployeePlacementHandler},
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

	return CreateBusinessResponse{
		BusinessID: business.ID.String(),
		Name:       business.Name,
		Code:       business.Code,
		Status:     business.Status,
		Plan:       business.Plan,
		Timezone:   business.Timezone,
	}, nil
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

func branchResponse(branch domain.Branch) BranchResponse {
	return BranchResponse{BranchID: branch.ID.String(), BusinessID: branch.BusinessID.String(), Name: branch.Name, Code: branch.Code, Address: branch.Address, Phone: branch.Phone, Status: branch.Status}
}
