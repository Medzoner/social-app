package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"social-app/api/auth"
	"social-app/internal/config"
	"social-app/internal/models"
	"social-app/internal/utils"
	"social-app/pkg/notifier"
	"social-app/pkg/ws"
)

type UseCase struct {
	mailer            notifier.Mailerx
	sms               notifier.SMSNotifier
	repo              Repository
	conn              *ws.Connector
	googleOAuthConfig *oauth2.Config
	jwtKey            []byte
	authCfg           config.Auth
}

func NewUseCase(r Repository, cfg config.Auth, c *ws.Connector, s notifier.SMSNotifier, m notifier.Mailerx) UseCase {
	googleOAuthConfig := &oauth2.Config{
		ClientID:     cfg.Sso.Google.ClientID,
		ClientSecret: cfg.Sso.Google.ClientSecret,
		RedirectURL:  cfg.Sso.Google.RedirectURL,
		Scopes:       cfg.Sso.Google.GetScopes(),
		Endpoint:     google.Endpoint,
	}
	return UseCase{
		repo:              r,
		jwtKey:            []byte(cfg.JWTSecret),
		conn:              c,
		mailer:            m,
		sms:               s,
		authCfg:           cfg,
		googleOAuthConfig: googleOAuthConfig,
	}
}

func (u UseCase) Register(ctx context.Context, input auth.RegisterInput) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	err := u.repo.Register(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}

	return nil
}

func (u UseCase) Login(ctx context.Context, input auth.LoginInput) (auth.JWTToken, error) {
	user, err := u.repo.Login(ctx, input.Username)
	if err != nil {
		return auth.JWTToken{}, fmt.Errorf("failed to login user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return auth.JWTToken{}, fmt.Errorf("invalid credentials: %w", err)
	}

	if !user.Verified || user.VerifiedExpires.Before(time.Now()) {
		if err := u.handleVerificationPending(ctx, user); err != nil {
			return auth.JWTToken{}, fmt.Errorf("failed to handle verification pending: %w", err)
		}
		return auth.JWTToken{
			VerifyResponse: auth.VerifyResponse{
				ID:       user.ID,
				Username: user.Username,
				Verified: false,
			},
		}, nil
	}

	token, err := u.generateJWT(user)
	if err != nil {
		return auth.JWTToken{}, fmt.Errorf("failed to generate JWT token: %w", err)
	}

	return token, nil
}

func (u UseCase) handleVerificationPending(ctx context.Context, user models.User) error {
	code, err := utils.GenerateCode(6)
	if err != nil {
		return fmt.Errorf("failed to generate verification code: %w", err)
	}
	user.VerifiedCode = code
	user.VerifiedAt = time.Now()
	user.VerifiedExpires = time.Now().Add(15 * 24 * time.Hour)

	if err := u.repo.Update2FA(ctx, user); err != nil {
		return fmt.Errorf("failed to save 2FA: %w", err)
	}

	if user.Phone != "" {
		if err := u.sms.SendPhoneVerification(user.Phone, code); err != nil {
			return fmt.Errorf("failed to send SMS verification: %w", err)
		}
	}
	if user.Email != "" {
		if err := u.mailer.SendEmailVerification(user.Email, code); err != nil {
			return fmt.Errorf("failed to send email verification: %w", err)
		}
	}

	return nil
}

func (u UseCase) Logout(ctx context.Context, userID uint64) error {
	if err := u.conn.CloseClient(ctx, userID); err != nil {
		return fmt.Errorf("failed to close websocket connection: %w", err)
	}

	return nil
}

func (u UseCase) RefreshToken(ctx context.Context, userID uint64) (auth.JWTToken, error) {
	refreshedUser, err := u.repo.Get(ctx, userID)
	if err != nil {
		return auth.JWTToken{}, fmt.Errorf("failed to refresh token: %w", err)
	}

	token, err := u.generateJWT(refreshedUser)
	if err != nil {
		return auth.JWTToken{}, fmt.Errorf("failed to generate JWT token: %w", err)
	}

	return token, nil
}

func (u UseCase) generateJWT(user models.User) (auth.JWTToken, error) {
	at, err := generateToken(user, u.authCfg.JWTExpDuration, map[string]any{}, u.jwtKey)
	if err != nil {
		return auth.JWTToken{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	rt, err := generateToken(user, u.authCfg.RefreshExpDuration, map[string]any{
		"typ": "refresh",
	}, u.jwtKey)
	if err != nil {
		return auth.JWTToken{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	id, err := generateToken(user, u.authCfg.JWTExpDuration, map[string]any{
		"role":  user.Role,
		"email": user.Email,
	}, u.jwtKey)
	if err != nil {
		return auth.JWTToken{}, fmt.Errorf("failed to generate ID token: %w", err)
	}

	return auth.JWTToken{
		AccessToken:  at,
		RefreshToken: rt,
		IDToken:      id,
	}, nil
}

func generateToken(user models.User, ttl time.Duration, extra map[string]any, key []byte) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"verified": user.Verified,
		"exp":      now.Add(ttl).Unix(),
		"iat":      now.Unix(),
		"nbf":      now.Unix(),
		"iss":      "social-app",
		"aud":      "social-app",
		"jti":      fmt.Sprintf("%d-%d", user.ID, now.Unix()),
	}

	for k, v := range extra {
		claims[k] = v
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ts, err := tok.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return ts, nil
}

func (u UseCase) Verify(ctx context.Context, userID uint64, code string) error {
	user, err := u.repo.GetProfile(ctx, fmt.Sprint(userID))
	if err != nil {
		return fmt.Errorf("failed to get user profile: %w", err)
	}

	if user.VerifiedCode != code {
		return fmt.Errorf("invalid code")
	}

	user.Verified = true
	user.VerifiedCode = ""
	user.VerifiedValidateAt = time.Now().Add(30 * 24 * time.Hour)

	return u.repo.Update(ctx, user)
}

func (u UseCase) SendEmailVerification(ctx context.Context, userID uint64) error {
	user, err := u.repo.GetProfile(ctx, fmt.Sprint(userID))
	if err != nil {
		return fmt.Errorf("failed to get user profile: %w", err)
	}

	code, err := utils.GenerateCode(6)
	if err != nil {
		return fmt.Errorf("failed to generate verification code: %w", err)
	}
	user.VerifiedCode = code

	if err := u.mailer.SendEmailVerification(user.Email, code); err != nil {
		return fmt.Errorf("failed to send email verification: %w", err)
	}

	return u.repo.Update(ctx, user)
}

func (u UseCase) SendPhoneVerification(ctx context.Context, userID uint64) error {
	user, err := u.repo.GetProfile(ctx, fmt.Sprint(userID))
	if err != nil {
		return fmt.Errorf("failed to get user profile: %w", err)
	}

	code, err := utils.GenerateCode(6)
	if err != nil {
		return fmt.Errorf("failed to generate verification code: %w", err)
	}
	user.VerifiedCode = code

	if err := u.sms.SendPhoneVerification(user.Phone, code); err != nil {
		return fmt.Errorf("failed to send phone verification: %w", err)
	}

	return u.repo.Update(ctx, user)
}

func (u UseCase) IsUserOnline(ctx context.Context, id uint64) (bool, error) {
	exists, err := u.conn.IsUserOnline(ctx, id)
	if err != nil {
		return false, fmt.Errorf("failed to check if user is online: %w", err)
	}
	return exists, nil
}

func (u UseCase) OauthCallback(ctx context.Context, input auth.OauthInput) (auth.JWTToken, error) {
	token, err := u.googleOAuthConfig.Exchange(ctx, input.Code)
	if err != nil {
		return auth.JWTToken{}, fmt.Errorf("failed to exchange OAuth code: %w", err)
	}

	client := u.googleOAuthConfig.Client(ctx, token)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.googleapis.com/oauth2/v2/userinfo", http.NoBody)
	if err != nil {
		return auth.JWTToken{}, fmt.Errorf("creating request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return auth.JWTToken{}, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	var userInfo struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return auth.JWTToken{}, fmt.Errorf("failed to decode user info: %w", err)
	}

	user, err := u.repo.GetByEmail(ctx, userInfo.Email)
	if err != nil {
		return auth.JWTToken{}, fmt.Errorf("failed to login user: %w", err)
	}

	if !user.Verified || user.VerifiedExpires.Before(time.Now()) {
		if err := u.handleVerificationPending(ctx, user); err != nil {
			return auth.JWTToken{}, fmt.Errorf("failed to handle verification pending: %w", err)
		}
		return auth.JWTToken{
			VerifyResponse: auth.VerifyResponse{
				ID:       user.ID,
				Username: user.Username,
				Verified: false,
			},
		}, nil
	}

	if user.IsZero() {
		user = models.User{
			Username: userInfo.Name,
			Email:    userInfo.Email,
			Avatar:   userInfo.Picture,
			Role:     "user",
			Verified: true,
		}

		if err := u.repo.Register(ctx, user); err != nil {
			return auth.JWTToken{}, fmt.Errorf("failed to register user: %w", err)
		}
	} else {
		user.Username = userInfo.Name
		user.Avatar = userInfo.Picture
		if err := u.repo.Update(ctx, user); err != nil {
			return auth.JWTToken{}, fmt.Errorf("failed to update user: %w", err)
		}
	}

	tk, err := u.generateJWT(user)
	if err != nil {
		return auth.JWTToken{}, fmt.Errorf("failed to generate JWT token: %w", err)
	}

	return tk, nil
}

func (u UseCase) OauthLogin(ctx context.Context) string {
	return u.googleOAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}
