package config

import "time"

// ConfigProvider Interface
type ConfigProvider interface {
	GetJWTSecret() string
	GetAccessTokenExpiry() time.Duration
	GetIssuer() string
}

// Config Implementation
func (c *Config) GetJWTSecret() string {
	return c.JWT.JWTSecret
}

func (c *Config) GetAccessTokenExpiry() time.Duration {
	return c.JWT.AccessTokenExpired
}

func (c *Config) GetIssuer() string {
	return c.JWT.Issuer
}
