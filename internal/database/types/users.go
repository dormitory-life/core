package types

type UserRole = string

const (
	UserAdminRole   UserRole = "admin"
	UserStudentRole UserRole = "student"
)

type GetUsersRoleRequest struct {
	UserId string
}

type GetUsersRoleResponse struct {
	Role UserRole
}
