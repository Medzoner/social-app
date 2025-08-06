package models

type Like struct {
	Model
	User   User   `json:"user"`
	UserID uint64 `gorm:"index;uniqueIndex:idx_user_post" json:"user_id"`
	PostID uint64 `gorm:"index;uniqueIndex:idx_user_post" json:"post_id"`
}

type LikeStats struct {
	PostID             uint64 `json:"post_id"`
	TotalLikes         uint64 `json:"total_likes"`
	LikedByCurrentUser bool   `json:"liked_by_current_user"`
}
