package storage

import (
	"io"
	"time"

	"github.com/dormitory-life/core/internal/constants"
)

type GetEntityFilesRequest struct {
	Category    constants.FileCategory
	EntityId    string
	SubEntityId string
	Amount      int
}

type GetEntityFilesResponse struct {
	FilesInfo []FileInfo
}

type UploadRequest struct {
	Category    constants.FileCategory
	EntityId    string
	SubEntityId string
	PhotoId     string
	FileName    string
	Reader      io.Reader
	Size        int64
	MimeType    string
}

type UploadResult struct {
	URL      string `json:"url"`
	FilePath string `json:"file_path"`
	FileName string `json:"file_name"`
	Size     int64  `json:"size"`
}

type GetFileResult struct {
	URL      string `json:"url"`
	FilePath string `json:"file_path"`
	FileName string `json:"file_name"`
	Size     int64  `json:"size"`
}

type FileInfo struct {
	Path         string    `json:"path"`
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"last_modified"`
	URL          string    `json:"url"`
}

type DeleteFileRequest struct {
	Path        *string
	Category    constants.FileCategory
	EntityId    string
	SubEntityId string
	FileName    string
}

type DeleteAllRequest struct {
	Category    constants.FileCategory
	EntityId    string
	SubEntityId string
}
