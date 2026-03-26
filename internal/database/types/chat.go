package types

import (
	"time"
)

type ChatMessage struct {
	ID          string
	DormitoryID string
	UserID      string
	Email       string
	Text        string
	CreatedAt   time.Time
}

type (
	GetChatMessagesRequest struct {
		DormitoryID string
		Page        uint64
	}

	GetChatMessagesResponse struct {
		Messages []ChatMessage
	}
)

type (
	CreateChatMessageRequest struct {
		DormitoryID string
		UserID      string
		Text        string
	}

	CreateChatMessageResponse struct {
		ID string
	}
)
