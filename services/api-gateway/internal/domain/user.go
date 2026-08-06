package domain

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

type UpdateUserRequest struct {
	FullName string `json:"full_name"`
	Status   string `json:"status"`
}

type AssignRoleRequest struct {
	Role string `json:"role"`
}

type CreatePlacementRequest struct {
	BranchID       string `json:"branch_id"`
	Position       string `json:"position"`
	EmploymentType string `json:"employment_type"`
}

type UserResponse struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Status   string `json:"status"`
}

type ListUsersResponse struct {
	Users []UserResponse `json:"users"`
}

type PlacementResponse struct {
	PlacementID    string `json:"placement_id"`
	UserID         string `json:"user_id"`
	BusinessID     string `json:"business_id"`
	BranchID       string `json:"branch_id"`
	Position       string `json:"position"`
	EmploymentType string `json:"employment_type"`
	Status         string `json:"status"`
}
