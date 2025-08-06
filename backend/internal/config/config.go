package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	Mailer   Mailgun `envPrefix:"MAILER_"`
	SMS      SMS     `envPrefix:"SMS_"`
	DB       DB      `envPrefix:"DB_"`
	Auth     Auth
	LLM      LLM      `envPrefix:"LLM_"`
	Mailtrap Mailtrap `envPrefix:"MAILER_"`
	Redis    Redis    `envPrefix:"REDIS_"`
}

type DB struct {
	DSN string `env:"DSN" envDefault:"postgres://user:pass@localhost:5464/socialdb"`
}

type Auth struct {
	JWTSecret          string        `env:"JWT_SECRET"  envDefault:"secret"`
	JWTExp             string        `env:"JWT_EXP"     envDefault:"15m"`
	RefreshExp         string        `env:"REFRESH_EXP" envDefault:"24h"`
	Sso                Sso           `envPrefix:"SSO_"`
	JWTExpDuration     time.Duration `env:"-"`
	RefreshExpDuration time.Duration `env:"-"`
}

type Sso struct {
	Google SsoGoogle `envPrefix:"GOOGLE_"`
}

type SsoGoogle struct {
	ClientID     string          `env:"CLIENT_ID"     envDefault:"xxxxxx.apps.googleusercontent.com"`
	ClientSecret string          `env:"CLIENT_SECRET" envDefault:"GOCSPX-xxxxxxx"`
	RedirectURL  string          `env:"REDIRECT_URL"  envDefault:"http://localhost:3222/oauth/google/callback"`
	ScopesRaw    string          `env:"SCOPES_RAW"    envDefault:"https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile openid"`
	Endpoint     oauth2.Endpoint `env:"ENDPOINT"`
	Scopes       []string
}

func (g *SsoGoogle) GetScopes() []string {
	return strings.Fields(g.ScopesRaw)
}

type Redis struct {
	Host     string `env:"HOST"      envDefault:"localhost"`
	Password string `env:"PASSWORD"  envDefault:""`
	Port     int    `env:"PORT"      envDefault:"6379"`
	DB       int    `env:"DB"        envDefault:"0"`
	PoolSize int    `env:"POOL_SIZE" envDefault:"10"`
	MinIdle  int    `env:"MIN_IDLE"  envDefault:"5"`
}

type Mailtrap struct {
	Host     string `env:"HOST"      envDefault:"sandbox.smtp.mailtrap.io"`
	AuthUser string `env:"AUTH_USER" envDefault:"xxxxxxx"`
	AuthPass string `env:"AUTH_PASS" envDefault:"xxxxxxx"`
	From     string `env:"FROM"      envDefault:"medzux@gmail.com"`
	Port     int    `env:"PORT"      envDefault:"2525"`
}

type Mailgun struct {
	Host   string `env:"HOST"    envDefault:"smtp.mailtrap.io"`
	ApiKey string `env:"API_KEY"`
	From   string `env:"FROM"    envDefault:"medzux@gmail.com"`
}

type SMS struct {
	AccountID string `env:"ACCOUNT_ID" envDefault:"xxxxxxx"`
	AuthToken string `env:"AUTH_TOKEN" envDefault:"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"`
	From      string `env:"FROM"       envDefault:"+15005550006"`
}

type LLM struct {
	URL string `env:"LLM_URL" envDefault:"http://medzoner-srv.lan:11434/api"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	var err error
	if cfg.Auth.JWTExpDuration, err = time.ParseDuration(cfg.Auth.JWTExp); err != nil {
		return nil, fmt.Errorf("invalid JWT_EXP: %w", err)
	}

	if cfg.Auth.RefreshExpDuration, err = time.ParseDuration(cfg.Auth.RefreshExp); err != nil {
		return nil, fmt.Errorf("invalid REFRESH_EXP: %w", err)
	}

	cfg.Auth.Sso.Google.Endpoint = google.Endpoint

	return cfg, nil
}
