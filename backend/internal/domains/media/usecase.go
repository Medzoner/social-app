package media

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"social-app/internal/models"
)

type MimeInfo struct {
	Extension string
	Folder    string
}

var allowedMIMEs = map[string]MimeInfo{
	"image/jpeg": {".jpg", "image"},
	"image/png":  {".png", "image"},
	"image/gif":  {".gif", "image"},
	"image/webp": {".webp", "image"},

	"video/mp4":        {".mp4", "video"},
	"video/quicktime":  {".mov", "video"},
	"video/x-msvideo":  {".avi", "video"},
	"video/x-matroska": {".mkv", "video"},

	"audio/mpeg": {".mp3", "audio"},
	"audio/wav":  {".wav", "audio"},
	"audio/ogg":  {".ogg", "audio"},
	"audio/webm": {".webm", "audio"},
	"audio/mp4":  {".m4a", "audio"},
}

type UseCase struct {
	repo Repository
}

func NewUseCase(r Repository) UseCase {
	return UseCase{
		repo: r,
	}
}

func (u UseCase) GetMedia(ctx context.Context, mediaUuid uuid.UUID) (models.Media, error) {
	m, err := u.repo.GetMedia(ctx, mediaUuid)
	if err != nil {
		return models.Media{}, fmt.Errorf("failed to get media: %w", err)
	}
	return m, nil
}

func (u UseCase) GetMedias(ctx context.Context, mediaUuids []uuid.UUID) ([]models.Media, error) {
	m, err := u.repo.GetMedias(ctx, mediaUuids)
	if err != nil {
		return nil, fmt.Errorf("failed to get media: %w", err)
	}
	return m, nil
}

func (u UseCase) UploadImage(ctx context.Context, file *multipart.FileHeader, userID uint64) (models.Media, error) {
	mime := file.Header.Get("Content-Type")

	info, ok := allowedMIMEs[mime]
	if !ok {
		return models.Media{}, fmt.Errorf("unsupported file type: %s", mime)
	}

	newUuid := uuid.New()
	fileName := newUuid.String() + info.Extension

	dirPath := filepath.Join("uploads", info.Folder)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return models.Media{}, fmt.Errorf("failed to create upload directory: %w", err)
	}

	savePath := filepath.Join(dirPath, fileName)

	if file.Size < 0 {
		return models.Media{}, fmt.Errorf("invalid file size: %d", file.Size)
	}
	fileSize := file.Size

	md := models.Media{
		Model: models.Model{
			UUID: uuid.NullUUID{
				UUID:  newUuid,
				Valid: true,
			},
		},
		FileName: fileName,
		FileExt:  info.Extension,
		FileSize: uint64(fileSize),
		FileType: mime,
		FilePath: savePath,
		UserID:   userID,
	}

	m, err := u.repo.UploadImage(ctx, md)
	if err != nil {
		return models.Media{}, fmt.Errorf("failed to upload image metadata: %w", err)
	}

	if err := u.SaveUploadedFile(ctx, file, savePath); err != nil {
		return models.Media{}, fmt.Errorf("failed to save uploaded file: %w", err)
	}

	return m, nil
}

func (u UseCase) SaveUploadedFile(ctx context.Context, file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0o750); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}
