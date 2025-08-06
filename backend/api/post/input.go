package post

import "github.com/google/uuid"

type CreatePostInput struct {
	Content    string          `json:"content"`
	MediaUUIDs []uuid.NullUUID `json:"media_uuids" swaggerignore:"true"`
	UserID     uint64          `json:"user_id"`
}
