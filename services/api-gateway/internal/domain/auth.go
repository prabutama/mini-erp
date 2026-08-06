package domain

type Role string

const (
	RolePlatformAdmin Role = "Platform Admin"
	RoleBusinessAdmin Role = "Business Admin"
	RoleManager       Role = "Manager"
	RoleStaff         Role = "Staff"
)

type SignupRequest struct {
	BusinessName     string `json:"business_name"`
	BusinessTimezone string `json:"business_timezone"`
	AdminFullName    string `json:"admin_full_name"`
	AdminEmail       string `json:"admin_email"`
	AdminPassword    string `json:"admin_password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthContext struct {
	UserID            string   `json:"user_id"`
	BusinessID        string   `json:"business_id,omitempty"`
	Role              Role     `json:"role"`
	Permissions       []string `json:"permissions"`
	AssignedBranchIDs []string `json:"assigned_branch_ids,omitempty"`
	RequestID         string   `json:"request_id"`
}

type AuthSession struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         AuthContext `json:"user"`
}
