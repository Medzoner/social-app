package profile

import (
	"context"
	"fmt"

	"social-app/internal/connector"
	"social-app/internal/models"
	"gorm.io/gorm"
	"errors"
)

type ProfileRepository interface {
	GetProfile(ctx context.Context, id string) (models.User, error)
	Update(ctx context.Context, m models.User) error
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
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, nil
		}
		return models.User{}, fmt.Errorf("failed to get user with id %s: %w", id, result.Error)
	}
	return user, nil
}

func (r Repository) Update(ctx context.Context, user models.User) error {
	return r.conn.DB.Save(&user).Error
}
