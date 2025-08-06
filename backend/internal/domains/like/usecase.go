package like

import (
	"context"
	"fmt"

	"social-app/api/like"
	"social-app/internal/models"
)

type UseCase struct {
	repo Repository
}

func NewUseCase(r Repository) UseCase {
	return UseCase{
		repo: r,
	}
}

func (u UseCase) LikePost(ctx context.Context, like like.PostInput) (models.Like, error) {
	l := models.Like{
		UserID: like.UserID,
		PostID: like.PostID,
	}

	l, err := u.repo.LikePost(ctx, l)
	if err != nil {
		return models.Like{}, err
	}

	return l, err
}

func (u UseCase) UnlikePost(ctx context.Context, like like.PostInput) (models.Like, error) {
	l := models.Like{
		UserID: like.UserID,
		PostID: like.PostID,
	}

	l, err := u.repo.UnlikePost(ctx, l)
	if err != nil {
		return models.Like{}, err
	}

	return l, nil
}

func (u UseCase) GetLikes(ctx context.Context, postID uint64) ([]models.Like, error) {
	l, err := u.repo.GetLikes(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("error getting likes: %w", err)
	}

	return l, nil
}

func (u UseCase) GetStats(ctx context.Context, postIDs []uint64, userID uint64) ([]models.LikeStats, error) {
	stats, err := u.repo.GetStats(ctx, postIDs, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting like stats: %w", err)
	}

	return stats, nil
}

func (u UseCase) GetLikeStats(ctx context.Context, postIDs []uint64, userID uint64) ([]models.LikeStats, error) {
	stats, err := u.repo.GetLikeStats(ctx, userID, postIDs)
	if err != nil {
		return nil, fmt.Errorf("error getting like stats: %w", err)
	}

	return stats, nil
}
