package domain

type RoleName string

const (
	RolePlatformAdmin RoleName = "Platform Admin"
	RoleBusinessAdmin RoleName = "Business Admin"
	RoleManager       RoleName = "Manager"
	RoleStaff         RoleName = "Staff"
)
