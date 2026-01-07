package requestmodels

import "mime/multipart"

type (
	CreateDormitoryPhotosRequest struct {
		DormitoryId       string
		PhotoFilesHeaders []*multipart.FileHeader
	}

	CreateDormitoryPhotosResponse struct {
		CreatePhotoResponses []CreatePhotoResponse
	}

	CreatePhotoResponse struct {
		URL      string
		FilePath string
		FileName string
		Size     int64
	}
)

type (
	DeleteDormitoryPhotosRequest struct {
		DormitoryId string
	}

	DeleteDormitoryPhotosResponse struct {
		DormitoryId string
	}
)
