package models

import "time"

type User struct {
	VerifiedExpires    time.Time   `gorm:"null"          json:"-"`
	VerifiedValidateAt time.Time   `gorm:"null"          json:"-"`
	VerifiedAt         time.Time   `gorm:"null"          json:"-"`
	Username           string      `gorm:"unique"        json:"username"`
	Role               string      `json:"role"`
	Email              string      `gorm:"null"          json:"-"`
	Phone              string      `gorm:"null"          json:"-"`
	VerifiedCode       string      `gorm:"null"          json:"-"`
	Avatar             string      `json:"avatar"`
	Password           string      `json:"-"`
	Bio                string      `json:"bio"`
	SignupType         string      `gorm:"default:email" json:"signup_type"`
	AvatarMedia        AvatarMedia `gorm:"-"             json:"avatar_media,omitempty"`
	Model
	ID       uint64 `gorm:"primaryKey"    json:"id"`
	Verified bool   `gorm:"default:false" json:"verified"`
	Active   bool   `gorm:"default:true"  json:"active"`
}

func (u User) IsZero() bool {
	return u == User{}
}

type AvatarMedia struct {
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
	FileType string `json:"file_type"`
	FileExt  string `json:"file_ext"`
	FileSize uint64 `json:"file_size"`
}
