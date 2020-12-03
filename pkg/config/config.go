package config

import "os"

const (
	configKey = "config"
)

type key string

// Config - config structure.
type Config struct {
	auth0       Auth0
	emailDomain string
	setPassword string
}

// Auth0 - Auth0 config structure.
type Auth0 struct {
	domain       string
	clientID     string
	clientSecret string
}

// LoadFromEnv - load configuration from environment variables
func (c *Config) LoadFromEnv() {
	c.auth0.domain = os.Getenv("AUTH0_DOMAIN")
	c.auth0.clientID = os.Getenv("AUTH0_CLIENT_ID")
	c.auth0.clientSecret = os.Getenv("AUTH0_CLIENT_SECRET")
	c.emailDomain = os.Getenv("EMAIL_DOMAIN")
	c.setPassword = os.Getenv("SET_PASSWORD")
}

// FromEnv - create config instance from environment variables
func FromEnv() *Config {
	cfg := Config{
		auth0: Auth0{
			domain:       os.Getenv("AUTH0_DOMAIN"),
			clientID:     os.Getenv("AUTH0_CLIENT_ID"),
			clientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		},
		emailDomain: os.Getenv("EMAIL_DOMAIN"),
		setPassword: os.Getenv("SET_PASSWORD"),
	}
	return &cfg
}

// Key -- context key for config.
func Key() interface{} {
	var k = key(configKey)
	return k
}
