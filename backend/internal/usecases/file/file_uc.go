package file

import (
	"astro-planner/backend/internal/entities"
	"astro-planner/backend/internal/repository"
	"astro-planner/backend/internal/storage"
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/google/uuid"
)

type FileUseCase struct {
	Repo *repository.FileRepo
	S3   *storage.S3Storage
}

func NewFileUseCase(repo *repository.FileRepo, s3 *storage.S3Storage) *FileUseCase {
	return &FileUseCase{Repo: repo, S3: s3}
}

func (uc *FileUseCase) UploadFile(ctx context.Context, userID string, file multipart.File, header *multipart.FileHeader) (*entities.File, string, error) {
	ext := filepath.Ext(header.Filename)
	objectKey := fmt.Sprintf("users/%s/%s%s", userID, uuid.New().String(), ext)
	file.Seek(0, 0)
	contentType := header.Header.Get("ContentType")
	if err := uc.S3.Upload(ctx, objectKey, file, contentType); err != nil {
		return nil, "", fmt.Errorf("s3 upload error: &w", err)
	}
	f := &entities.File{
		UserID:       &userID,
		ObjectKey:    objectKey,
		OriginalName: header.Filename,
		ContentType:  contentType,
		Size:         header.Size,
	}
	if err := uc.Repo.Save(ctx, f); err != nil {
		return nil, "", fmt.Errorf("db save error: %w", err)
	}
	publicURL := uc.S3.GetPublicURL(objectKey)
	return f, publicURL, nil
}
