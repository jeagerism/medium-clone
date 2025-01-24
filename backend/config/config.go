package config

import (
	"fmt"
	"time"
)

// Config Struct
type Config struct {
	server *server
	db     *db
	jwt    *jwt
}

// jwt Struct (Implementation of ConfigProvider)
type jwt struct {
	SecretKey           string
	AccessTokenExpired  time.Duration
	RefreshTokenExpired time.Duration
	Issuer              string
}

// ConfigProvider Interface
type ConfigJWT interface {
	GetJWTSecret() []byte
	GetAccessTokenExpiry() time.Duration
	GetRefreshTokenExpiry() time.Duration
	GetIssuer() string
}

func (c *Config) JWT() ConfigJWT {
	return c.jwt
}

// Implementing ConfigProvider
func (j *jwt) GetJWTSecret() []byte {
	return []byte(j.SecretKey)
}

func (j *jwt) GetAccessTokenExpiry() time.Duration {
	return j.AccessTokenExpired
}

func (j *jwt) GetRefreshTokenExpiry() time.Duration {
	return j.RefreshTokenExpired
}

func (j *jwt) GetIssuer() string {
	return j.Issuer
}

type db struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
	sslMode  string
	timezone string
}

type ConfigDB interface {
	Url() string
}

func (c *Config) Db() ConfigDB {
	return c.db
}

func (db *db) Url() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		db.host,
		db.user,
		db.password,
		db.dbname,
		db.port,
		db.sslMode,
		db.timezone,
	)
}

type server struct {
	port int
}

type ConfigServer interface {
	GetPort() int
}

func (c *Config) Server() ConfigServer {
	return c.server
}

func (s *server) GetPort() int {
	return s.port
}
