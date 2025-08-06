package profile

import (
	"context"

	"social-app/internal/models"
	"social-app/pkg/notifier"
)

type UseCase struct {
	repo   Repository
	mailer notifier.Mailerx
	sms    notifier.SMSNotifier
}

func NewUseCase(r Repository, s notifier.SMSNotifier, m notifier.Mailerx) UseCase {
	return UseCase{
		repo:   r,
		mailer: m,
		sms:    s,
	}
}

func (u UseCase) GetProfile(ctx context.Context, id string) (models.User, error) {
	profile, err := u.repo.GetProfile(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	return profile, nil
}

func (u UseCase) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	err := u.repo.Update(ctx, user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
