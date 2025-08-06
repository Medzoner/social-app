package media

import (
	"context"

	"github.com/google/uuid"
	"social-app/internal/connector"
	"social-app/internal/models"
)

type MediaRepository interface {
	UploadImage(ctx context.Context, m models.Media) (models.Media, error)
	GetMedia(ctx context.Context, mediaUuid uuid.UUID) (models.Media, error)
	GetMedias(ctx context.Context, mediaUuids []uuid.UUID) ([]models.Media, error)
}

type Repository struct {
	conn connector.DBConn
}

func NewRepository(con connector.DBConn) Repository {
	return Repository{
		conn: con,
	}
}

func (r Repository) UploadImage(ctx context.Context, m models.Media) (models.Media, error) {
	result := r.conn.DB.Create(&m)
	if result.Error != nil {
		return models.Media{}, result.Error
	}

	return m, nil
}

func (r Repository) GetMedia(ctx context.Context, mediaUuid uuid.UUID) (models.Media, error) {
	var media models.Media
	result := r.conn.DB.First(&media, "uuid = ?", mediaUuid)
	if result.Error != nil {
		return models.Media{}, result.Error
	}

	return media, nil
}

func (r Repository) GetMedias(ctx context.Context, mediaUuids []uuid.UUID) ([]models.Media, error) {
	var medias []models.Media
	result := r.conn.DB.Where("uuid IN ?", mediaUuids).Find(&medias)
	if result.Error != nil {
		return nil, result.Error
	}

	return medias, nil
}
