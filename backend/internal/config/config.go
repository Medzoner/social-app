package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Mailer   Mailgun `envPrefix:"MAILER_"`
	SMS      SMS     `envPrefix:"SMS_"`
	DB       DB      `envPrefix:"DB_"`
	Auth     Auth
	Mailtrap Mailtrap `envPrefix:"MAILER_"`
	LLM      LLM      `envPrefix:"LLM_"`
}

type DB struct {
	DSN string `env:"DSN" envDefault:"postgres://user:pass@localhost:5464/socialdb"`
}

type Auth struct {
	JWTSecret string `env:"JWT_SECRET" envDefault:"secret"`
}

type Mailtrap struct {
	Host     string `env:"HOST"      envDefault:"sandbox.smtp.mailtrap.io"`
	AuthUser string `env:"AUTH_USER" envDefault:"89ef83a7cddb56"`
	AuthPass string `env:"AUTH_PASS" envDefault:"a09074a9f61630"`
	From     string `env:"FROM"      envDefault:"medzux@gmail.com"`
	Port     int    `env:"PORT"      envDefault:"2525"`
}

type Mailgun struct {
	Host   string `env:"HOST"    envDefault:"smtp.mailtrap.io"`
	ApiKey string `env:"API_KEY"`
	From   string `env:"FROM"    envDefault:"medzux@gmail.com"`
}

type SMS struct {
	AccountID string `env:"ACCOUNT_ID" envDefault:"xxx"`
	AuthToken string `env:"AUTH_TOKEN" envDefault:"xxx"`
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
	return cfg, nil
}
