package models

import "time"

type Comment struct {
	Content string `json:"content"`
	UserModel
	PostID uint64 `json:"post_id"`
}

func (c Comment) GetCursorFields() (createdAt time.Time, id uint64) {
	return c.CreatedAt, c.ID
}

type CommentList struct {
	NextCursor string    `json:"next_cursor,omitempty"`
	Comments   []Comment `json:"comments"`
	HasMore    bool      `json:"has_more"`
}

type CommentCount struct {
	Total  int64  `json:"total"`
	PostID uint64 `json:"post_id"`
}
