package comment

import (
	"context"
	"fmt"

	"social-app/api/comment"
	"social-app/internal/models"
	"social-app/pkg/ws"
)

type UseCase struct {
	bc   ws.Broadcaster
	repo Repository
}

func NewUseCase(r Repository, bc ws.Broadcaster) UseCase {
	return UseCase{
		repo: r,
		bc:   bc,
	}
}

func (u UseCase) CreateComment(ctx context.Context, c comment.CreateCommentInput) (models.Comment, error) {
	cm := models.Comment{
		PostID: c.PostID,
		UserModel: models.UserModel{
			UserID: c.UserID,
		},
		Content: c.Content,
	}

	result, err := u.repo.CreateComment(ctx, cm)
	if err != nil {
		return models.Comment{}, fmt.Errorf("error creating comment: %w", err)
	}

	u.bc.NotifyAll("New comment added!")

	return result, nil
}

func (u UseCase) GetComments(ctx context.Context, cursor string, postID uint64) (models.CommentList, error) {
	comments, err := u.repo.GetComments(ctx, cursor, postID)
	if err != nil {
		return models.CommentList{}, fmt.Errorf("error fetching comments: %w", err)
	}

	return comments, nil
}

func (u UseCase) CountByPost(ctx context.Context, postID uint64) (int64, error) {
	total, err := u.repo.CountByPost(ctx, postID)
	if err != nil {
		return 0, fmt.Errorf("error counting comments: %w", err)
	}
	return total, nil
}

func (u UseCase) CountByPosts(ctx context.Context, postIDs []uint64, userID uint64) ([]models.CommentCount, error) {
	cc, err := u.repo.CountByPosts(ctx, postIDs, userID)
	if err != nil {
		return []models.CommentCount{}, fmt.Errorf("error counting comments: %w", err)
	}
	return cc, nil
}
