package models

import (
	"time"
)

type Post struct {
	Content  string    `json:"content"`
	Comments []Comment `json:"comments"`
	Likes    []Like    `json:"likes"`
	Medias   []Media   `gorm:"many2many:post_medias;" json:"medias"`
	Model
	User   User   `gorm:"foreignKey:UserID;references:ID" json:"user"`
	UserID uint64 `gorm:"index"                           json:"user_id"`
}

func (p Post) GetCursorFields() (createdAt time.Time, id uint64) {
	return p.CreatedAt, p.ID
}

type PostList struct {
	NextCursor string `json:"next_cursor,omitempty"`
	Posts      []Post `json:"posts"`
	HasMore    bool   `json:"has_more"`
}
