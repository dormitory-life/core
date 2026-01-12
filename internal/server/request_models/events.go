package requestmodels

import (
	"fmt"
	"mime/multipart"
	"net/url"
	"time"

	dbtypes "github.com/dormitory-life/core/internal/database/types"
)

type Event struct {
	EventId     string     `json:"event_id"`
	DormitoryId string     `json:"dormitory_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	EventPhotos []FileInfo `json:"event_photos"`
	CreatedAt   time.Time  `json:"created_at"`
}

type (
	GetDormitoryEventsRequest struct {
		DormitoryId string
		Page        uint64
	}

	GetDormitoryEventsResponse struct {
		Events []Event `json:"events"`
	}
)

func (*GetDormitoryEventsRequest) FromUrlQuery(query url.Values) (*GetDormitoryEventsRequest, error) {
	res := &GetDormitoryEventsRequest{
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

func (r *Event) From(msg *dbtypes.Event) *Event {
	if msg == nil {
		return nil
	}

	return &Event{
		EventId:     msg.EventId,
		DormitoryId: msg.DormitoryId,
		Title:       msg.Title,
		Description: msg.Description,
		CreatedAt:   msg.CreatedAt,
	}
}

func (*GetDormitoryEventsResponse) From(msg *dbtypes.GetDormitoryEventsResponse) *GetDormitoryEventsResponse {
	if msg == nil {
		return nil
	}

	res := &GetDormitoryEventsResponse{
		Events: make([]Event, 0),
	}

	for _, val := range msg.Events {
		review := *new(Event).From(&val)
		res.Events = append(res.Events, review)
	}

	return res
}

type (
	CreateDormitoryEventRequest struct {
		DormitoryId       string
		PhotoFilesHeaders []*multipart.FileHeader
		Title             string
		Description       string
	}

	CreateDormitoryEventResponse struct {
		EventId              string                `json:"event_id"`
		DormitoryId          string                `json:"dormitory_id"`
		CreatePhotoResponses []CreatePhotoResponse `json:"photos"`
		Title                string                `json:"title"`
		Description          string                `json:"description"`
	}
)

type (
	DeleteDormitoryEventRequest struct {
		DormitoryId string
		EventId     string
	}

	DeleteDormitoryEventResponse struct {
	}
)
