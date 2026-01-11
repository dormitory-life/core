package types

import "time"

type Review struct {
	ReviewId    string
	OwnerId     string
	DormitoryId string
	Title       string
	Description string
	CreatedAt   time.Time
}

type (
	GetDormitoryReviewsRequest struct {
		DormitoryId string
		Page        uint64
	}

	GetDormitoryReviewsResponse struct {
		Reviews []Review
	}
)

type (
	CreateReviewRequest struct {
		OwnerId     string
		DormitoryId string
		Title       string
		Description string
	}

	CreateReviewResponse struct {
		ReviewId string
	}
)

type (
	DeleteReviewRequest struct {
		ReviewId string
	}

	DeleteReviewResponse struct {
	}
)
