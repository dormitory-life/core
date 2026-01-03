package requestmodels

import (
	dbtypes "github.com/dormitory-life/core/internal/database/types"
)

type Dormitory struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Address       string `json:"address"`
	Support_email string `json:"support_email"`
	Description   string `json:"description"`
}

type (
	GetDormitoriesRequest  struct{}
	GetDormitoriesResponse struct {
		Dormitories []Dormitory `json:"dormitories"`
	}
)

func (r *Dormitory) From(msg *dbtypes.Dormitory) *Dormitory {
	if msg == nil {
		return nil
	}

	return &Dormitory{
		Id:            msg.Id,
		Name:          msg.Name,
		Address:       msg.Address,
		Support_email: msg.Support_email,
		Description:   msg.Description,
	}
}

func (r *GetDormitoriesResponse) From(msg *dbtypes.GetDormitoriesResponse) *GetDormitoriesResponse {
	if msg == nil {
		return nil
	}

	res := &GetDormitoriesResponse{
		Dormitories: make([]Dormitory, 0),
	}

	for _, val := range msg.Dormitories {
		dormitory := *new(Dormitory).From(&val)
		res.Dormitories = append(res.Dormitories, dormitory)
	}

	return res
}

type (
	GetDormitoryByIdRequest struct {
		DormitoryId string `json:"dormitory_id"`
	}

	GetDormitoryByIdResponse struct {
		Dormitory Dormitory `json:"dormitory"`
	}
)

func (r *GetDormitoryByIdResponse) From(msg *dbtypes.GetDormitoryByIdResponse) *GetDormitoryByIdResponse {
	if msg == nil {
		return nil
	}

	return &GetDormitoryByIdResponse{
		Dormitory: Dormitory{
			Id:            msg.Dormitory.Id,
			Name:          msg.Dormitory.Name,
			Address:       msg.Dormitory.Address,
			Support_email: msg.Dormitory.Support_email,
			Description:   msg.Dormitory.Description,
		},
	}
}
