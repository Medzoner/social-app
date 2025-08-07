package comment

type CreateCommentInput struct {
	Content string `binding:"required" json:"content"`
	PostID  uint64 `binding:"required" json:"post_id"`
	UserID  uint64 `json:"user_id"`
}
