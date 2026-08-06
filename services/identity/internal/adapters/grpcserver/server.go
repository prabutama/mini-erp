package grpcserver

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/isapr/mini-erp/services/identity/internal/application"
	"github.com/isapr/mini-erp/services/identity/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/status"
)

const serviceName = "identity.v1.IdentityService"

type SignupTenantAdminRequest struct {
	BusinessID    string `json:"business_id"`
	AdminEmail    string `json:"admin_email"`
	AdminPassword string `json:"admin_password"`
	AdminFullName string `json:"admin_full_name"`
	RequestID     string `json:"request_id"`
}

type SignupTenantAdminResponse struct {
	AccessToken       string   `json:"access_token"`
	RefreshToken      string   `json:"refresh_token"`
	UserID            string   `json:"user_id"`
	BusinessID        string   `json:"business_id"`
	Role              string   `json:"role"`
	Permissions       []string `json:"permissions"`
	AssignedBranchIDs []string `json:"assigned_branch_ids"`
	RequestID         string   `json:"request_id"`
}

type LoginRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	RequestID string `json:"request_id"`
}

type LoginResponse = SignupTenantAdminResponse

type GetUserAccessContextRequest struct {
	AccessToken string `json:"access_token"`
	RequestID   string `json:"request_id"`
}

type GetUserAccessContextResponse struct {
	UserID            string   `json:"user_id"`
	BusinessID        string   `json:"business_id,omitempty"`
	Role              string   `json:"role"`
	Permissions       []string `json:"permissions"`
	AssignedBranchIDs []string `json:"assigned_branch_ids,omitempty"`
	RequestID         string   `json:"request_id"`
}

type UserResponse struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Status   string `json:"status"`
}

type CreateUserRequest struct {
	BusinessID string `json:"business_id"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	FullName   string `json:"full_name"`
	Role       string `json:"role"`
}

type ListUsersRequest struct {
	BusinessID string `json:"business_id"`
}

type ListUsersResponse struct {
	Users []UserResponse `json:"users"`
}

type GetUserRequest struct {
	UserID string `json:"user_id"`
}

type UpdateUserRequest struct {
	UserID   string `json:"user_id"`
	FullName string `json:"full_name"`
	Status   string `json:"status"`
}

type AssignBusinessRoleRequest struct {
	UserID     string `json:"user_id"`
	BusinessID string `json:"business_id"`
	Role       string `json:"role"`
}

type AssignBusinessRoleResponse struct {
	UserID     string `json:"user_id"`
	BusinessID string `json:"business_id"`
	Role       string `json:"role"`
}

type Server struct {
	signup *application.SignupService
	auth   *application.AuthService
	users  *application.UserService
}

type identityServiceServer interface {
	SignupTenantAdmin(context.Context, SignupTenantAdminRequest) (SignupTenantAdminResponse, error)
	Login(context.Context, LoginRequest) (LoginResponse, error)
	GetUserAccessContext(context.Context, GetUserAccessContextRequest) (GetUserAccessContextResponse, error)
	CreateUser(context.Context, CreateUserRequest) (UserResponse, error)
	ListUsers(context.Context, ListUsersRequest) (ListUsersResponse, error)
	GetUser(context.Context, GetUserRequest) (UserResponse, error)
	UpdateUser(context.Context, UpdateUserRequest) (UserResponse, error)
	AssignBusinessRole(context.Context, AssignBusinessRoleRequest) (AssignBusinessRoleResponse, error)
}

func New(signup *application.SignupService, auth *application.AuthService, users *application.UserService) *Server {
	encoding.RegisterCodec(jsonCodec{})
	return &Server{signup: signup, auth: auth, users: users}
}

func (s *Server) Register(grpcServer *grpc.Server) {
	grpcServer.RegisterService(&grpc.ServiceDesc{
		ServiceName: serviceName,
		HandlerType: (*identityServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			{MethodName: "SignupTenantAdmin", Handler: signupTenantAdminHandler},
			{MethodName: "Login", Handler: loginHandler},
			{MethodName: "GetUserAccessContext", Handler: getUserAccessContextHandler},
			{MethodName: "CreateUser", Handler: createUserHandler},
			{MethodName: "ListUsers", Handler: listUsersHandler},
			{MethodName: "GetUser", Handler: getUserHandler},
			{MethodName: "UpdateUser", Handler: updateUserHandler},
			{MethodName: "AssignBusinessRole", Handler: assignBusinessRoleHandler},
		},
	}, s)
}

func signupTenantAdminHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req SignupTenantAdminRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(identityServiceServer).SignupTenantAdmin(ctx, req)
}

func loginHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req LoginRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(identityServiceServer).Login(ctx, req)
}

func getUserAccessContextHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req GetUserAccessContextRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(identityServiceServer).GetUserAccessContext(ctx, req)
}

func createUserHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req CreateUserRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(identityServiceServer).CreateUser(ctx, req)
}

func listUsersHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req ListUsersRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(identityServiceServer).ListUsers(ctx, req)
}

func getUserHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req GetUserRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(identityServiceServer).GetUser(ctx, req)
}

func updateUserHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req UpdateUserRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(identityServiceServer).UpdateUser(ctx, req)
}

func assignBusinessRoleHandler(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
	var req AssignBusinessRoleRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	return srv.(identityServiceServer).AssignBusinessRole(ctx, req)
}

func (s *Server) SignupTenantAdmin(ctx context.Context, req SignupTenantAdminRequest) (SignupTenantAdminResponse, error) {
	businessID, err := uuid.Parse(req.BusinessID)
	if err != nil {
		return SignupTenantAdminResponse{}, status.Error(codes.InvalidArgument, "invalid business_id")
	}

	session, err := s.signup.SignupTenantAdminSession(ctx, application.SignupTenantAdminInput{
		BusinessID: businessID,
		Email:      req.AdminEmail,
		Password:   req.AdminPassword,
		FullName:   req.AdminFullName,
	}, req.RequestID)
	if errors.Is(err, application.ErrValidation) {
		return SignupTenantAdminResponse{}, status.Error(codes.InvalidArgument, "validation failed")
	}
	if errors.Is(err, application.ErrUserAlreadyExists) {
		return SignupTenantAdminResponse{}, status.Error(codes.AlreadyExists, "user already exists")
	}
	if err != nil {
		return SignupTenantAdminResponse{}, status.Error(codes.Internal, err.Error())
	}

	log.Printf("grpc service=identity method=SignupTenantAdmin user_id=%s business_id=%s role=%s", session.Context.UserID, session.Context.BusinessID, session.Context.Role)

	return sessionResponse(session), nil
}

func (s *Server) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	session, err := s.auth.Login(ctx, application.LoginInput{Email: req.Email, Password: req.Password, RequestID: req.RequestID})
	if errors.Is(err, application.ErrValidation) {
		return LoginResponse{}, status.Error(codes.InvalidArgument, "validation failed")
	}
	if errors.Is(err, application.ErrInvalidCredentials) {
		return LoginResponse{}, status.Error(codes.Unauthenticated, "invalid credentials")
	}
	if err != nil {
		return LoginResponse{}, status.Error(codes.Internal, err.Error())
	}
	return sessionResponse(session), nil
}

func (s *Server) GetUserAccessContext(ctx context.Context, req GetUserAccessContextRequest) (GetUserAccessContextResponse, error) {
	authContext, err := s.auth.GetUserAccessContext(ctx, req.AccessToken, req.RequestID)
	if errors.Is(err, application.ErrInvalidToken) {
		return GetUserAccessContextResponse{}, status.Error(codes.Unauthenticated, "invalid token")
	}
	if err != nil {
		return GetUserAccessContextResponse{}, status.Error(codes.Internal, err.Error())
	}
	return contextResponse(authContext), nil
}

func (s *Server) CreateUser(ctx context.Context, req CreateUserRequest) (UserResponse, error) {
	businessID, err := uuid.Parse(req.BusinessID)
	if err != nil {
		return UserResponse{}, status.Error(codes.InvalidArgument, "invalid business_id")
	}
	user, err := s.users.CreateBusinessUser(ctx, application.CreateUserInput{BusinessID: businessID, Email: req.Email, Password: req.Password, FullName: req.FullName, Role: domain.RoleName(req.Role)})
	if errors.Is(err, application.ErrValidation) {
		return UserResponse{}, status.Error(codes.InvalidArgument, "validation failed")
	}
	if errors.Is(err, application.ErrUserAlreadyExists) {
		return UserResponse{}, status.Error(codes.AlreadyExists, "user already exists")
	}
	if err != nil {
		return UserResponse{}, status.Error(codes.Internal, err.Error())
	}
	return userResponse(user), nil
}

func (s *Server) ListUsers(ctx context.Context, req ListUsersRequest) (ListUsersResponse, error) {
	businessID, err := uuid.Parse(req.BusinessID)
	if err != nil {
		return ListUsersResponse{}, status.Error(codes.InvalidArgument, "invalid business_id")
	}
	users, err := s.users.ListBusinessUsers(ctx, businessID)
	if err != nil {
		return ListUsersResponse{}, status.Error(codes.Internal, err.Error())
	}
	response := ListUsersResponse{Users: []UserResponse{}}
	for _, user := range users {
		response.Users = append(response.Users, userResponse(user))
	}
	return response, nil
}

func (s *Server) GetUser(ctx context.Context, req GetUserRequest) (UserResponse, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return UserResponse{}, status.Error(codes.InvalidArgument, "invalid user_id")
	}
	user, err := s.users.GetUser(ctx, userID)
	if err != nil {
		return UserResponse{}, status.Error(codes.NotFound, "user not found")
	}
	return userResponse(user), nil
}

func (s *Server) UpdateUser(ctx context.Context, req UpdateUserRequest) (UserResponse, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return UserResponse{}, status.Error(codes.InvalidArgument, "invalid user_id")
	}
	user, err := s.users.UpdateUser(ctx, application.UpdateUserInput{UserID: userID, FullName: req.FullName, Status: req.Status})
	if err != nil {
		return UserResponse{}, status.Error(codes.Internal, err.Error())
	}
	return userResponse(user), nil
}

func (s *Server) AssignBusinessRole(ctx context.Context, req AssignBusinessRoleRequest) (AssignBusinessRoleResponse, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return AssignBusinessRoleResponse{}, status.Error(codes.InvalidArgument, "invalid user_id")
	}
	businessID, err := uuid.Parse(req.BusinessID)
	if err != nil {
		return AssignBusinessRoleResponse{}, status.Error(codes.InvalidArgument, "invalid business_id")
	}
	if err := s.users.AssignBusinessRole(ctx, userID, domain.RoleName(req.Role), businessID); err != nil {
		return AssignBusinessRoleResponse{}, status.Error(codes.Internal, err.Error())
	}
	return AssignBusinessRoleResponse{UserID: req.UserID, BusinessID: req.BusinessID, Role: req.Role}, nil
}

func sessionResponse(session domain.AuthSession) SignupTenantAdminResponse {
	ctx := contextResponse(session.Context)
	return SignupTenantAdminResponse{
		AccessToken:       session.AccessToken,
		RefreshToken:      session.RefreshToken,
		UserID:            ctx.UserID,
		BusinessID:        ctx.BusinessID,
		Role:              ctx.Role,
		Permissions:       ctx.Permissions,
		AssignedBranchIDs: ctx.AssignedBranchIDs,
		RequestID:         ctx.RequestID,
	}
}

func contextResponse(authContext domain.AuthContext) GetUserAccessContextResponse {
	resp := GetUserAccessContextResponse{
		UserID:            authContext.UserID.String(),
		Role:              string(authContext.Role),
		Permissions:       authContext.Permissions,
		AssignedBranchIDs: authContext.AssignedBranchIDs,
		RequestID:         authContext.RequestID,
	}
	if authContext.BusinessID != uuid.Nil {
		resp.BusinessID = authContext.BusinessID.String()
	}
	return resp
}

func userResponse(user domain.User) UserResponse {
	return UserResponse{UserID: user.ID.String(), Email: user.Email, FullName: user.FullName, Status: user.Status}
}
