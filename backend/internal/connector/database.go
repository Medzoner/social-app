package connector

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"social-app/internal/config"
	"social-app/internal/models"
)

type DBConn struct {
	DB *gorm.DB
}

func NewDBConn(cfg config.DB) (DBConn, error) {
	db, err := ConnectDatabase(cfg.DSN)
	if err != nil {
		return DBConn{}, fmt.Errorf("failed to connect to database: %w", err)
	}
	return DBConn{DB: db}, nil
}

func ConnectDatabase(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)

	if err := db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Comment{},
		&models.Like{},
		&models.Message{},
		&models.Notification{},
		&models.Media{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
