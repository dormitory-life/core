package emailer

type SendMessageRequest struct {
	UserEmail    string
	SupportEmail string
	Title        string
	Description  string
}

type SendMessageResponse struct {
	UserEmail    string
	SupportEmail string
}
