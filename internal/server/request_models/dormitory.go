package requestmodels

import (
	"time"

	dbtypes "github.com/dormitory-life/core/internal/database/types"
	"github.com/dormitory-life/core/internal/storage"
)

type Dormitory struct {
	Id           string     `json:"id"`
	Name         string     `json:"name"`
	Address      string     `json:"address"`
	SupportEmail string     `json:"support_email"`
	Description  string     `json:"description"`
	Photos       []FileInfo `json:"photo_links"`
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
		Id:           msg.Id,
		Name:         msg.Name,
		Address:      msg.Address,
		SupportEmail: msg.SupportEmail,
		Description:  msg.Description,
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
			Id:           msg.Dormitory.Id,
			Name:         msg.Dormitory.Name,
			Address:      msg.Dormitory.Address,
			SupportEmail: msg.Dormitory.SupportEmail,
			Description:  msg.Dormitory.Description,
		},
	}
}

type FileInfo struct {
	Path         string    `json:"path"`
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"last_modified"`
	URL          string    `json:"url"`
}

func ConvertFileInfos(msg *storage.GetEntityFilesResponse) []FileInfo {
	fileInfos := make([]FileInfo, 0)
	for _, info := range msg.FilesInfo {
		fileInfos = append(fileInfos, FileInfo{
			Path:         info.Path,
			Name:         info.Name,
			Size:         info.Size,
			LastModified: info.LastModified,
			URL:          info.URL,
		})
	}

	return fileInfos
}

type (
	CreateDormitoryRequest struct {
		DormitoryId  string `json:"dormitory_id"`
		Name         string `json:"name"`
		Address      string `json:"address"`
		SupportEmail string `json:"support_email"`
		Description  string `json:"description"`
	}

	CreateDormitoryResponse struct {
		DormitoryId string `json:"dormitory_id"`
	}
)

func (r *CreateDormitoryResponse) From(msg *dbtypes.CreateDormitoryResponse) *CreateDormitoryResponse {
	if msg == nil {
		return nil
	}

	return &CreateDormitoryResponse{
		DormitoryId: msg.DormitoryId,
	}
}

type (
	UpdateDormitoryRequest struct {
		DormitoryId  string
		Name         *string `json:"name"`
		Address      *string `json:"address"`
		SupportEmail *string `json:"support_email"`
		Description  *string `json:"description"`
	}

	UpdateDormitoryResponse struct {
		DormitoryId string `json:"dormitory_id"`
	}
)

func (r *UpdateDormitoryResponse) From(msg *dbtypes.UpdateDormitoryResponse) *UpdateDormitoryResponse {
	if msg == nil {
		return nil
	}

	return &UpdateDormitoryResponse{
		DormitoryId: msg.DormitoryId,
	}
}

type (
	DeleteDormitoryRequest struct {
		DormitoryId string `json:"dormitory_id"`
	}

	DeleteDormitoryResponse struct{}
)
