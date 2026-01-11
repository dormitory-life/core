package requestmodels

import (
	"fmt"
	"mime/multipart"
	"net/url"
	"strconv"
	"time"

	"github.com/dormitory-life/core/internal/constants"
	dbtypes "github.com/dormitory-life/core/internal/database/types"
)

type Review struct {
	ReviewId     string     `json:"review_id"`
	OwnerId      string     `json:"owner_id"`
	DormitoryId  string     `json:"dormitory_id"`
	ReviewPhotos []FileInfo `json:"review_photos"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	CreatedAt    time.Time  `json:"created_at"`
}

type (
	GetDormitoryReviewsRequest struct {
		DormitoryId string
		Page        uint64
	}

	GetDormitoryReviewsResponse struct {
		Reviews []Review `json:"reviews"`
	}
)

func (*GetDormitoryReviewsRequest) FromUrlQuery(query url.Values) (*GetDormitoryReviewsRequest, error) {
	res := &GetDormitoryReviewsRequest{
		Page: constants.DefaultReviewsPageSize,
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

func (r *Review) From(msg *dbtypes.Review) *Review {
	if msg == nil {
		return nil
	}

	return &Review{
		ReviewId:    msg.ReviewId,
		OwnerId:     msg.OwnerId,
		DormitoryId: msg.DormitoryId,
		Title:       msg.Title,
		Description: msg.Description,
		CreatedAt:   msg.CreatedAt,
	}
}

func (*GetDormitoryReviewsResponse) From(msg *dbtypes.GetDormitoryReviewsResponse) *GetDormitoryReviewsResponse {
	if msg == nil {
		return nil
	}

	res := &GetDormitoryReviewsResponse{
		Reviews: make([]Review, 0),
	}

	for _, val := range msg.Reviews {
		review := *new(Review).From(&val)
		res.Reviews = append(res.Reviews, review)
	}

	return res
}

type (
	CreateReviewRequest struct {
		OwnerId           string
		DormitoryId       string
		PhotoFilesHeaders []*multipart.FileHeader
		Title             string
		Description       string
	}

	CreateReviewResponse struct {
		ReviewId             string                `json:"review_id"`
		OwnerId              string                `json:"owner_id"`
		DormitoryId          string                `json:"dormitory_id"`
		CreatePhotoResponses []CreatePhotoResponse `json:"photos"`
		Title                string                `json:"title"`
		Description          string                `json:"description"`
	}
)

type (
	DeleteReviewRequest struct {
		DormitoryId string
		ReviewId    string
	}

	DeleteReviewResponse struct {
	}
)

func parseUint64(val string) (uint64, error) {
	uintVal, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, err
	}

	return uintVal, nil
}
