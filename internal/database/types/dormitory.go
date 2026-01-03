package types

type Dormitory struct {
	Id           string
	Name         string
	Address      string
	SupportEmail string
	Description  string
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

type (
	CreateDormitoryRequest struct {
		DormitoryId  string
		Name         string
		Address      string
		SupportEmail string
		Description  string
	}

	CreateDormitoryResponse struct {
		DormitoryId string
	}
)

type (
	UpdateDormitoryRequest struct {
		DormitoryId  string
		Name         *string
		Address      *string
		SupportEmail *string
		Description  *string
	}

	UpdateDormitoryResponse struct {
		DormitoryId string
	}
)
