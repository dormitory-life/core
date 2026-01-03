package types

type Dormitory struct {
	Id            string
	Name          string
	Address       string
	Support_email string
	Description   string
}

type (
	GetDormitoriesRequest  struct{}
	GetDormitoriesResponse struct {
		Dormitories []Dormitory
	}
)

type (
	GetDormitoryByIdRequest struct {
		DormitoryId string
	}

	GetDormitoryByIdResponse struct {
		Dormitory Dormitory
	}
)
