package constants

const (
	PathDormitoryPhotos = "dormitory/%s/photos/"
	PathReviewPhotos    = "dormitory/%s/reviews/%s/photos/"
	PathFeedPhotos      = "dormitory/%s/feed/%s/photos/"
)

type FileCategory = string

const (
	CategoryDormitoryPhotos FileCategory = "dormitory"
	CategoryEventPhotos     FileCategory = "event"
	CategoryReviewPhotos    FileCategory = "review"
)

const (
	GetDormitoriesDefaultAmount int = 1
)
