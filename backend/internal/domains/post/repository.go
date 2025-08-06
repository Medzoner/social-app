package post

import (
	"context"
	"fmt"

	"social-app/internal/connector"
	"social-app/internal/models"
	"social-app/pkg/pagination"
)

const limit = 10

type PostRepository interface {
	CreatePost(ctx context.Context, p models.Post) (models.Post, error)
	GetPosts(ctx context.Context, cursor, ownerId string) (models.PostList, error)
}

type Repository struct {
	conn connector.DBConn
}

func NewRepository(con connector.DBConn) Repository {
	return Repository{
		conn: con,
	}
}

func (r Repository) CreatePost(ctx context.Context, post models.Post) (models.Post, error) {
	result := r.conn.DB.Create(&post)
	if result.Error != nil {
		return models.Post{}, result.Error
	}

	return post, nil
}

func (r Repository) GetPosts(ctx context.Context, cursor, ownerId string) (models.PostList, error) {
	var posts []models.Post

	db := r.conn.DB.
		Preload("User").
		Preload("Medias").
		Limit(limit + 1)

	if ownerId != "" {
		db = db.
			Where("user_id = ?", ownerId)
	}

	db = pagination.CursorFilter[models.Post](cursor, db, "desc")

	if err := db.Find(&posts).Error; err != nil {
		return models.PostList{}, fmt.Errorf("failed to fetch posts: %w", err)
	}

	nextCursor, hasMore, items := pagination.NextCursor(posts, limit)

	return models.PostList{
		Posts:      items,
		HasMore:    hasMore,
		NextCursor: nextCursor,
	}, nil
}
