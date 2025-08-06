package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserModel struct {
	Model
	User     User      `json:"user"`
	UserID   uint64    `gorm:"index" json:"user_id"   param:"user_id"`
	UserUUID uuid.UUID `gorm:"index" json:"user_uuid"`
}

type Model struct {
	CreatedAt time.Time      `gorm:"index"                                json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                                json:"deleted_at"                         swaggerignore:"true"`
	ID        uint64         `example:"123"                               gorm:"primaryKey;index:idx_id_deleted_at" json:"id"`
	UUID      uuid.NullUUID  `gorm:"type:uuid;default:uuid_generate_v4()" json:"uuid"                               swaggerignore:"true"`
}
