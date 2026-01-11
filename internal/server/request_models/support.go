package requestmodels

type CreateSupportRequest struct {
	UserId      string
	DormitoryId string
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateSupportResponse struct {
	Message string `json:"message"`
}
