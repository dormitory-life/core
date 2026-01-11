package types

type (
	GetEmailsForSupportRequest struct {
		UserId      string
		DormitoryId string
	}

	GetEmailsForSupportResponse struct {
		UserEmail    string
		SupportEmail string
	}
)
