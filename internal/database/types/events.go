package types

import "time"

type Event struct {
	EventId     string
	DormitoryId string
	Title       string
	Description string
	CreatedAt   time.Time
}

type (
	GetDormitoryEventsRequest struct {
		DormitoryId string
		Page        uint64
	}

	GetDormitoryEventsResponse struct {
		Events []Event
	}
)

type (
	CreateDormitoryEventRequest struct {
		DormitoryId string
		Title       string
		Description string
	}

	CreateDormitoryEventResponse struct {
		EventId     string
		DormitoryId string
		Title       string
		Description string
	}
)

type (
	DeleteDormitoryEventRequest struct {
		EventId string
	}

	DeleteDormitoryEventResponse struct {
	}
)
