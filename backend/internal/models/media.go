package models

import (
	"gorm.io/gorm"
)

type Media struct {
	DeletedAt gorm.DeletedAt `gorm:"index"          json:"-"`
	FileName  string         `gorm:"size:255;index" json:"file_name"`
	FilePath  string         `gorm:"size:512"       json:"file_path"`
	FileType  string         `gorm:"size:50;index"  json:"file_type"`
	FileExt   string         `gorm:"size:10;index"  json:"file_ext"`
	Model
	User     User   `gorm:"foreignKey:UserID" json:"user"`
	UserID   uint64 `gorm:"index"             json:"user_id"`
	FileSize uint64 `json:"file_size"`
}
