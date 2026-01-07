package storage

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/dormitory-life/core/internal/constants"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3StorageConfig struct {
	Type            string
	Endpoint        string
	AccessKeyId     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
	PublicUrl       string

	Logger slog.Logger
}

type MinIOClient struct {
	client    *minio.Client
	logger    slog.Logger
	bucket    string
	publicUrl string
}

type Storage interface {
	// GetFile - получение файла по пути в S3 по относительному пути в S3
	GetFile(ctx context.Context, filePath string) (io.ReadCloser, error)

	// GetEntityFiles - получение списка файлов по путям /{category}/{entityId}/
	GetEntityFiles(ctx context.Context, category constants.FileCategory, entityId string) ([]FileInfo, error)

	// Upload загружает файл в папку /{category}/{entityId}/
	Upload(ctx context.Context, req *UploadRequest) (*UploadResult, error)

	// Delete удаляет файл по относительному пути в S3
	Delete(ctx context.Context, filePath string) error

	// DeleteAll удаляет все файлы по пути /{category}/{entityId}/
	DeleteAll(ctx context.Context, category constants.FileCategory, entityId string) error

	// GetFileUrl - вспомогательная функция по получению ссылки для пользователя на файл в S3
	GetFileURL(filePath string) string

	// GetMimeType - вспомогательная функция по получению типа контента по имени файла
	GetMimeType(filename string) string
}

func New(cfg S3StorageConfig) (Storage, error) {
	switch cfg.Type {
	case "minio":
		return newMinIOClient(cfg)
	default:
		return nil, fmt.Errorf("unknown storage type: %s", cfg.Type)
	}
}

func newMinIOClient(cfg S3StorageConfig) (*MinIOClient, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyId, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	cli := &MinIOClient{
		client:    client,
		logger:    cfg.Logger,
		bucket:    cfg.BucketName,
		publicUrl: cfg.PublicUrl,
	}

	if err := cli.initBucket(context.Background()); err != nil {
		return nil, err
	}

	return cli, nil
}

func (m *MinIOClient) initBucket(ctx context.Context) error {
	exists, err := m.client.BucketExists(ctx, m.bucket)
	if err != nil {
		return fmt.Errorf("failed to check bucket: %w", err)
	}

	if !exists {
		err = m.client.MakeBucket(ctx, m.bucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		fmt.Printf("Bucket '%s' created\n", m.bucket)
	}

	return nil
}

func (m *MinIOClient) GetFile(ctx context.Context, filePath string) (io.ReadCloser, error) {
	obj, err := m.client.GetObject(ctx, m.bucket, filePath, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	_, err = obj.Stat()
	if err != nil {
		obj.Close()
		return nil, fmt.Errorf("file not found: %w", err)
	}

	return obj, nil
}

func (m *MinIOClient) GetEntityFiles(ctx context.Context, category constants.FileCategory, entityId string) ([]FileInfo, error) {
	prefix := fmt.Sprintf("%s/%s/", category, entityId)

	var files []FileInfo

	objectsCh := m.client.ListObjects(ctx, m.bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for obj := range objectsCh {
		if obj.Err != nil {
			continue
		}

		files = append(files, FileInfo{
			Path: obj.Key,
			Name: extractFileName(obj.Key),
			Size: obj.Size,
			URL:  m.GetFileURL(obj.Key),
		})
	}

	return files, nil
}

func (m *MinIOClient) Upload(
	ctx context.Context,
	req *UploadRequest,
) (*UploadResult, error) {
	if req == nil {
		return nil, fmt.Errorf("upload request is nil")
	}

	m.logger.Debug("uploading file",
		slog.String("category", string(req.Category)),
		slog.String("entity", req.EntityId),
		slog.String("photo_id", req.PhotoId),
		slog.String("filename", req.FileName),
	)

	fileExt := getFileExtension(req.FileName)
	fileName := fmt.Sprintf("%s_%s%s",
		getFilePrefix(req.Category),
		req.PhotoId,
		fileExt,
	)

	s3Path := fmt.Sprintf("%s/%s/%s", req.Category, req.EntityId, fileName)

	_, err := m.client.PutObject(ctx, m.bucket, s3Path, req.Reader, req.Size, minio.PutObjectOptions{
		ContentType: req.MimeType,
	})
	if err != nil {
		m.logger.Error("failed to upload file",
			slog.String("error", err.Error()),
			slog.String("path", s3Path),
		)
		return nil, fmt.Errorf("failed to upload file to S3: %w", err)
	}

	publicURL := m.GetFileURL(s3Path)

	m.logger.Debug("uploaded photo", slog.String("file_path", s3Path), slog.String("file_name", fileName))

	return &UploadResult{
		URL:      publicURL,
		FilePath: s3Path,
		FileName: fileName,
		Size:     req.Size,
	}, nil
}

func (m *MinIOClient) Delete(
	ctx context.Context,
	filePath string,
) error {
	err := m.client.RemoveObject(ctx, m.bucket, filePath, minio.RemoveObjectOptions{})
	if err != nil {
		m.logger.Error("failed to delete file",
			slog.String("error", err.Error()),
			slog.String("path", filePath),
		)
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// DeleteAll - удаляем все файлы для сущности по префиксу
func (m *MinIOClient) DeleteAll(
	ctx context.Context,
	category constants.FileCategory,
	entityId string,
) error {
	prefix := fmt.Sprintf("%s/%s/", category, entityId)

	objectsCh := m.client.ListObjects(ctx, m.bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for obj := range objectsCh {
		if obj.Err != nil {
			m.logger.Error("error listing objects", slog.String("error", obj.Err.Error()))
			continue
		}

		err := m.Delete(ctx, obj.Key)
		if err != nil {
			m.logger.Warn("failed to delete object",
				slog.String("object", obj.Key),
				slog.String("error", err.Error()),
			)
		}
	}

	return nil
}

func (m *MinIOClient) GetFileURL(filePath string) string {
	return fmt.Sprintf("%s/%s", m.publicUrl, filePath)
}

func extractFileName(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return path
	}
	return parts[len(parts)-1]
}

func getFileExtension(filename string) string {
	dotIndex := strings.LastIndex(filename, ".")
	if dotIndex == -1 || dotIndex == len(filename)-1 {
		return ""
	}
	return filename[dotIndex:]
}

func getFilePrefix(category constants.FileCategory) string {
	switch category {
	case constants.CategoryDormitoryPhotos:
		return "DormPhoto"
	case constants.CategoryEventPhotos:
		return "EventPhoto"
	case constants.CategoryReviewPhotos:
		return "ReviewPhoto"
	default:
		return "File"
	}
}

func (m *MinIOClient) GetMimeType(filename string) string {
	ext := strings.ToLower(getFileExtension(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}
