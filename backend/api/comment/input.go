package comment

type CreateCommentInput struct {
	Content string `json:"content"`
	PostID  uint64 `json:"post_id"`
	UserID  uint64 `json:"user_id"`
}
