package post

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"social-app/api/post"
	"social-app/internal/domains/media"
	"social-app/internal/models"
	"social-app/pkg/ws"
)

type UseCase struct {
	bc   ws.Broadcaster
	repo Repository
	mUC  media.UseCase
}

func NewUseCase(r Repository, m media.UseCase, bc ws.Broadcaster) UseCase {
	return UseCase{
		repo: r,
		mUC:  m,
		bc:   bc,
	}
}

func (u UseCase) CreatePost(ctx context.Context, input post.CreatePostInput) (models.Post, error) {
	p := models.Post{
		Content: input.Content,
		UserID:  input.UserID,
	}

	mediaUuids := make([]uuid.UUID, 0, len(input.MediaUUIDs))
	for _, mu := range input.MediaUUIDs {
		if mu.Valid {
			mediaUuids = append(mediaUuids, mu.UUID)
		}
	}

	m, err := u.mUC.GetMedias(ctx, mediaUuids)
	if err != nil {
		return models.Post{}, fmt.Errorf("error fetching media: %w", err)
	}

	p.Medias = m

	result, err := u.repo.CreatePost(ctx, p)
	if err != nil {
		return p, fmt.Errorf("error creating post: %w", err)
	}

	u.bc.NotifyAll("{\"type\": \"post_created\", \"data\": \"new post\"}")
	return result, nil
}

func (u UseCase) GetPosts(ctx context.Context, cursor, id string) (models.PostList, error) {
	p, err := u.repo.GetPosts(ctx, cursor, id)
	if err != nil {
		return models.PostList{}, fmt.Errorf("error fetching posts: %w", err)
	}

	for i := range p.Posts {
		avatar, err := u.getAvatar(ctx, p.Posts[i].User)
		if err != nil {
			return models.PostList{}, fmt.Errorf("error fetching avatar for post %d: %w", p.Posts[i].ID, err)
		}
		p.Posts[i].User.AvatarMedia = avatar
	}

	return p, nil
}

func (u UseCase) getAvatar(ctx context.Context, user models.User) (models.AvatarMedia, error) {
	if user.Avatar == "" {
		return models.AvatarMedia{}, nil
	}
	m, err := u.mUC.GetMedia(ctx, uuid.MustParse(user.Avatar))
	if err != nil {
		return models.AvatarMedia{}, fmt.Errorf("failed to get media for user %d: %w", user.ID, err)
	}
	return models.AvatarMedia{
		FileName: m.FileName,
		FileSize: m.FileSize,
		FileType: m.FileType,
		FilePath: m.FilePath,
	}, nil
}
