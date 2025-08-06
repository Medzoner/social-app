package like

import (
	"context"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
	"social-app/internal/connector"
	"social-app/internal/models"
)

type LikeRepository interface {
	LikePost(ctx context.Context, l models.Like) (models.Like, error)
	UnlikePost(ctx context.Context, l models.Like) (models.Like, error)
	GetLikes(ctx context.Context, postID uint64) ([]models.Like, error)
	GetLikeStats(ctx context.Context, userID, postID uint64) (models.LikeStats, error)
}

type Repository struct {
	conn connector.DBConn
}

func NewRepository(con connector.DBConn) Repository {
	return Repository{
		conn: con,
	}
}

func (r Repository) LikePost(ctx context.Context, like models.Like) (models.Like, error) {
	var existing models.Like

	// On cherche même les likes soft-deleted
	err := r.conn.DB.Unscoped().
		Where("user_id = ? AND post_id = ?", like.UserID, like.PostID).
		First(&existing).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Like{}, fmt.Errorf("error checking like existence: %w", err)
	}

	// Cas 1: Le like existe (actif ou soft-deleted)
	if existing.ID != 0 {
		// Si supprimé, on le restaure
		if existing.DeletedAt.Valid {
			if err := r.conn.DB.
				Unscoped().
				Model(&existing).
				UpdateColumn("deleted_at", gorm.DeletedAt{}).Error; err != nil {
				return models.Like{}, fmt.Errorf("failed to restore like: %w", err)
			}
			log.Printf("✅ Like restored for user %d on post %d", like.UserID, like.PostID)
		}
		return existing, nil
	}

	// Cas 2: Le like n'existe pas encore
	if err := r.conn.DB.Create(&like).Error; err != nil {
		return models.Like{}, fmt.Errorf("error creating like: %w", err)
	}

	return like, nil
}

func (r Repository) UnlikePost(ctx context.Context, like models.Like) (models.Like, error) {
	if err := r.conn.DB.
		Where("user_id = ? AND post_id = ?", like.UserID, like.PostID).
		First(&like).Error; err != nil {
		return like, fmt.Errorf("like not found: %w", err)
	}

	if err := r.conn.DB.Delete(&like).Error; err != nil {
		return like, err
	}
	return like, nil
}

func (r Repository) GetLikes(ctx context.Context, postID uint64) ([]models.Like, error) {
	var likes []models.Like
	r.conn.DB.
		Where("post_id = ? AND deleted_at IS NULL", postID).
		Find(&likes)

	return likes, r.conn.DB.Error
}

func (r Repository) GetLikeStats(ctx context.Context, userID uint64, postIDs []uint64) ([]models.LikeStats, error) {
	var stats []models.LikeStats

	subQuery := r.conn.DB.
		Table("likes").
		Select("post_id, true AS liked_by_current_user").
		Where("user_id = ? AND post_id IN ?", userID, postIDs)

	err := r.conn.DB.
		Table("likes").
		Select(`
			likes.post_id,
			COUNT(likes.user_id) AS total_likes,
			COALESCE(user_likes.liked_by_current_user, false) AS liked_by_current_user
		`).
		Where("likes.post_id IN ?", postIDs).
		Joins("LEFT JOIN (?) AS user_likes ON likes.post_id = user_likes.post_id", subQuery).
		Group("likes.post_id, user_likes.liked_by_current_user").
		Scan(&stats).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch like stats: %w", err)
	}

	return stats, nil
}

func (r Repository) GetStats(ctx context.Context, ids []uint64, userID uint64) ([]models.LikeStats, error) {
	stats := make([]models.LikeStats, 0)

	if len(ids) == 0 {
		return stats, nil
	}

	query := `
		SELECT 
			post_id,
			COUNT(*) AS total_likes,
			COUNT(CASE WHEN user_id = ? THEN 1 END) > 0 AS liked_by_current_user
		FROM likes
		WHERE post_id IN (?) AND deleted_at IS NULL
		GROUP BY post_id
	`

	err := r.conn.DB.Raw(query, userID, ids).Scan(&stats).Error
	if err != nil {
		return nil, fmt.Errorf("error getting like stats: %w", err)
	}

	return stats, nil
}
