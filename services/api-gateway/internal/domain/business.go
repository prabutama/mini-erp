package domain

type BusinessResponse struct {
	BusinessID    string `json:"business_id"`
	Name          string `json:"name"`
	Code          string `json:"code"`
	Status        string `json:"status"`
	Plan          string `json:"plan"`
	PlatformNotes string `json:"platform_notes,omitempty"`
	SuspendedAt   string `json:"suspended_at,omitempty"`
	Timezone      string `json:"timezone"`
	CreatedAt     string `json:"created_at,omitempty"`
	UpdatedAt     string `json:"updated_at,omitempty"`
}

type ListBusinessesResponse struct {
	Businesses []BusinessResponse `json:"businesses"`
}

type UpdateBusinessRequest struct {
	Name     string `json:"name"`
	Timezone string `json:"timezone"`
}

type UpdatePlatformBusinessRequest struct {
	Status        string `json:"status"`
	Plan          string `json:"plan"`
	PlatformNotes string `json:"platform_notes"`
	SuspendedAt   string `json:"suspended_at"`
}

type RoleResponse struct {
	Name  Role   `json:"name"`
	Scope string `json:"scope"`
}

type ListRolesResponse struct {
	Roles []RoleResponse `json:"roles"`
}
