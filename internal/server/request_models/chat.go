package requestmodels

import (
	"fmt"
	"net/url"
	"time"

	dbtypes "github.com/dormitory-life/core/internal/database/types"
)

type ChatMessage struct {
	ID          string    `json:"id"`
	DormitoryID string    `json:"dormitory_id"`
	UserID      string    `json:"user_id"`
	UserEmail   string    `json:"email"`
	Text        string    `json:"text"`
	CreatedAt   time.Time `json:"created_at"`
}

type (
	GetChatMessagesRequest struct {
		DormitoryID string
		Page        uint64
	}

	GetChatMessagesResponse struct {
		Messages []ChatMessage `json:"messages"`
	}
)

func (*ChatMessage) From(msg *dbtypes.ChatMessage) *ChatMessage {
	if msg == nil {
		return nil
	}

	return &ChatMessage{
		ID:          msg.ID,
		DormitoryID: msg.DormitoryID,
		UserID:      msg.UserID,
		UserEmail:   msg.Email,
		Text:        msg.Text,
		CreatedAt:   msg.CreatedAt,
	}
}

func (*GetChatMessagesRequest) FromUrlQuery(query url.Values) (*GetChatMessagesRequest, error) {
	res := &GetChatMessagesRequest{
		Page: 1,
	}

	if query == nil {
		return res, nil
	}

	if val, ok := query["page"]; ok {
		intVal, err := parseUint64(val[0])
		if err != nil {
			return nil, fmt.Errorf("invalid page param: %w", err)
		}

		res.Page = intVal
	}

	if res.Page == 0 {
		res.Page = 1
	}

	return res, nil
}

func (*GetChatMessagesResponse) From(msg *dbtypes.GetChatMessagesResponse) *GetChatMessagesResponse {
	if msg == nil {
		return nil
	}

	res := &GetChatMessagesResponse{
		Messages: make([]ChatMessage, 0),
	}

	for _, val := range msg.Messages {
		message := *new(ChatMessage).From(&val)
		res.Messages = append(res.Messages, message)
	}

	return res
}

type (
	CreateChatMessageRequest struct {
		DormitoryID string `json:"dormitory_id"`
		Text        string `json:"text"`
	}

	CreateChatMessageResponse struct {
		ID string `json:"id"`
	}
)

func (*CreateChatMessageResponse) From(msg *dbtypes.CreateChatMessageResponse) *CreateChatMessageResponse {
	if msg == nil {
		return nil
	}

	return &CreateChatMessageResponse{
		ID: msg.ID,
	}
}
