package config

import (
	"context"
	"os"
)

const (
	configKey            = "config"
	envAuth0Domain       = "AUTH0_DOMAIN"
	envAuth0ClientID     = "AUTH0_CLIENT_ID"
	envAuth0ClientSecret = "AUTH0_CLIENT_SECRET" //nolint: gosec
	envEmailDomain       = "EMAIL_DOMAIN"
	envSetPassword       = "SET_PASSWORD"
)

type key string

// Config - config structure.
type Config struct {
	Auth0       Auth0
	EmailDomain string
	SetPassword string
}

// Auth0 - Auth0 config structure.
type Auth0 struct {
	Domain       string
	ClientID     string
	ClientSecret string
}

// FromEnv - create config instance from environment variables
func FromEnv() *Config {
	cfg := Config{
		Auth0: Auth0{
			Domain:       os.Getenv(envAuth0Domain),
			ClientID:     os.Getenv(envAuth0ClientID),
			ClientSecret: os.Getenv(envAuth0ClientSecret),
		},
		EmailDomain: os.Getenv(envEmailDomain),
		SetPassword: os.Getenv(envSetPassword),
	}
	return &cfg
}

// Key -- context key for config.
func Key() interface{} {
	var k = key(configKey)
	return k
}

// FromContext -- extract config from context value.
func FromContext(ctx context.Context) *Config {
	cfg, ok := ctx.Value(Key()).(*Config)
	if !ok {
		return nil
	}

	return cfg
}
