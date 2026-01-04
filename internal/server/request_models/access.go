package requestmodels

type CheckAccessRequest struct {
	UserId       string
	DormitoryId  string
	RoleRequired bool
}

type CheckAccessResponse struct {
	Allowed bool
	Reason  string
	Error   error
}
