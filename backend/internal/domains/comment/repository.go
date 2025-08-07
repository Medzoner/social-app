package comment

import (
	"context"
	"fmt"

	"social-app/internal/connector"
	"social-app/internal/models"
	"social-app/pkg/pagination"
)

const limit = 5

type CommentRepository interface {
	CreateComment(ctx context.Context, com models.Comment) (models.Comment, error)
	GetComments(ctx context.Context, cursor string, postID uint64) (models.CommentList, error)
	CountByPost(ctx context.Context, postID uint64) (int64, error)
}

type Repository struct {
	conn connector.DBConn
}

func NewRepository(con connector.DBConn) Repository {
	return Repository{
		conn: con,
	}
}

func (r Repository) CreateComment(ctx context.Context, cm models.Comment) (models.Comment, error) {
	var post models.Post
	if err := r.conn.DB.First(&post, cm.PostID).Error; err != nil {
		return models.Comment{}, fmt.Errorf("post not found: %w", err)
	}

	if result := r.conn.DB.
		Preload("User").
		Create(&cm); result.Error != nil {
		return cm, result.Error
	}

	return cm, nil
}

func (r Repository) GetComments(ctx context.Context, cursor string, postID uint64) (models.CommentList, error) {
	var comments []models.Comment

	db := r.conn.DB

	db = db.
		Preload("User").
		Where("post_id = ?", postID).
		Limit(limit + 1)

	db = pagination.CursorFilter[models.Comment](cursor, db, "desc")

	if err := db.Find(&comments).Error; err != nil {
		return models.CommentList{}, err
	}

	nextCursor, hasMore, items := pagination.NextCursor(comments, limit)

	return models.CommentList{
		NextCursor: nextCursor,
		Comments:   items,
		HasMore:    hasMore,
	}, nil
}

func (r Repository) CountByPost(ctx context.Context, postID uint64) (int64, error) {
	var total int64

	db := r.conn.DB

	db.Model(&models.Comment{}).
		Where("post_id = ?", postID).
		Count(&total)

	return total, nil
}

func (r Repository) CountByPosts(ctx context.Context, ids []uint64, _ uint64) ([]models.CommentCount, error) {
	ccs := make([]models.CommentCount, 0)

	if len(ids) == 0 {
		return ccs, nil
	}

	query := `
		SELECT 
			post_id,
			COUNT(id) AS total
		FROM comments
		WHERE post_id IN (?) AND deleted_at IS NULL
		GROUP BY post_id
	`

	err := r.conn.DB.Raw(query, ids).Scan(&ccs).Error
	if err != nil {
		return nil, fmt.Errorf("error getting like stats: %w", err)
	}

	return ccs, nil
}
