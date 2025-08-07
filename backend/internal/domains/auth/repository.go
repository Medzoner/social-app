package auth

import (
	"context"
	"fmt"

	"social-app/internal/connector"
	"social-app/internal/models"
	"gorm.io/gorm"
	"errors"
)

type AuthRepository interface {
	Register(ctx context.Context, user models.User) (models.User, error)
	Login(ctx context.Context, username string) (models.User, error)
	GetProfile(ctx context.Context, id string) (models.User, error)
	Update(ctx context.Context, user models.User) error
	Save(user models.User) error
	Get(ctx context.Context, id uint64) (models.User, error)
	Update2FA(ctx context.Context, user models.User) error
	GetByEmail(ctx context.Context, email string) (models.User, error)
}

type Repository struct {
	conn connector.DBConn
}

func NewRepository(con connector.DBConn) Repository {
	return Repository{
		conn: con,
	}
}

func (r Repository) GetProfile(ctx context.Context, id string) (models.User, error) {
	var user models.User
	result := r.conn.DB.
		First(&user, id)
	if result.Error != nil {
		return models.User{}, fmt.Errorf("failed to get user with id %s: %w", id, result.Error)
	}
	return user, nil
}

func (r Repository) Save(user models.User) error {
	result := r.conn.DB.
		Save(&user)
	if result.Error != nil {
		return fmt.Errorf("failed to save user: %w", result.Error)
	}
	return nil
}

func (r Repository) Update(ctx context.Context, user models.User) error {
	return r.conn.DB.Save(&user).Error
}

func (r Repository) Register(ctx context.Context, user models.User) (models.User, error) {
	result := r.conn.DB.
		Create(&user)
	if result.Error != nil {
		return user, fmt.Errorf("failed to create user: %w", result.Error)
	}
	return user, nil
}

func (r Repository) Login(ctx context.Context, username string) (models.User, error) {
	var user models.User
	result := r.conn.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return models.User{}, fmt.Errorf("failed to find user: %w", result.Error)
	}

	return user, nil
}

func (r Repository) Get(ctx context.Context, id uint64) (models.User, error) {
	var user models.User
	result := r.conn.DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return models.User{}, fmt.Errorf("failed to find user: %w", result.Error)
	}

	return user, nil
}

func (r Repository) Update2FA(ctx context.Context, user models.User) error {
	return r.conn.DB.Model(&user).
		Updates(map[string]interface{}{
			"verified_code":    user.VerifiedCode,
			"verified_expires": user.VerifiedExpires,
		}).Error
}

func (r Repository) GetByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	result := r.conn.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, nil // No user found with the given email
		}
		return user, fmt.Errorf("failed to find user by email: %w", result.Error)
	}

	return user, nil
}
