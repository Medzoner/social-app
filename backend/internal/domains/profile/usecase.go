package profile

import (
	"context"

	"social-app/internal/models"
	"social-app/pkg/notifier"
	"social-app/internal/domains/media"
	"net/url"
	"log"
	"github.com/google/uuid"
	"fmt"
)

type UseCase struct {
	repo    Repository
	mediaUC media.UseCase
	mailer  notifier.Mailerx
	sms     notifier.SMSNotifier
}

func NewUseCase(r Repository, muc media.UseCase, s notifier.SMSNotifier, m notifier.Mailerx) UseCase {
	return UseCase{
		repo:    r,
		mailer:  m,
		sms:     s,
		mediaUC: muc,
	}
}

func (u UseCase) GetProfile(ctx context.Context, id string) (models.User, error) {
	profile, err := u.repo.GetProfile(ctx, id)
	if err != nil {
		return profile, fmt.Errorf("failed to get profile for user %s: %w", id, err)
	}

	profile, err = u.getAvatar(ctx, profile)
	if err != nil {
		log.Printf("Failed to get avatar for user %s: %v", id, err)
		return profile, err
	}

	return profile, nil
}

func (u UseCase) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	err := u.repo.Update(ctx, user)
	if err != nil {
		return models.User{}, err
	}

	p, err := u.getAvatar(ctx, user)
	if err != nil {
		log.Printf("Failed to get avatar for user %d: %v", user.ID, err)
		return user, err
	}
	return p, nil
}

func (u UseCase) getAvatar(ctx context.Context, profile models.User) (models.User, error) {
	if profile.Avatar != "" {
		uri, err := url.ParseRequestURI(profile.Avatar)
		if err != nil {
			log.Printf("Invalid avatar URL for user %d: %v", profile.ID, err)
		}

		if uri != nil {
			profile.AvatarMedia = models.AvatarMedia{
				FileName: uri.Path,
				FilePath: uri.String(),
			}

			return profile, nil
		}

		m, err := u.mediaUC.GetMedia(ctx, uuid.MustParse(profile.Avatar))
		if err != nil {
			log.Printf("Failed to get media for user %d: %v", profile.ID, err)
			return profile, err
		}

		profile.AvatarMedia = models.AvatarMedia{
			FileName: m.FileName,
			FileSize: m.FileSize,
			FileType: m.FileType,
			FilePath: u.mediaUC.ResolvePath(m.FilePath),
		}
	}
	return profile, nil
}
