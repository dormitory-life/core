package constants

const (
	SchemaName string = "public"
)

const (
	UsersTableName              string = "users"
	DormitoryTableName          string = "dormitory"
	GradesTable                 string = "grades"
	DormitoryAvgGradesTableName string = "dormitory_average_grades"
	ReviewTableName             string = "reviews"
	FeedTableName               string = "feed"
)

const (
	DefaultReviewsPageSize    uint64 = 10
	DefaultPaginationPageSize uint64 = 10
	DefaultEventsPageSize     uint64 = 15
)
