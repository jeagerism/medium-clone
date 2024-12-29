package config

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Server *Server
		DB     *DB
		JWT    *JWT // Add JWT here to ensure it's included in the final config structure
	}

	Server struct {
		Port int
	}

	DB struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		SSLMode  string
		Timezone string
	}

	JWT struct {
		JWTSecret           string
		AccessTokenExpired  time.Duration
		RefreshTokenExpired time.Duration
		Issuer              string
	}
)

var (
	once           sync.Once
	configInstance *Config
)

func GetConfig() *Config {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./config")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file: %v", err)
		}

		// Unmarshal the general config into Config instance
		if err := viper.Unmarshal(&configInstance); err != nil {
			log.Fatalf("Unable to decode into struct: %v", err)
		}

		// Parse duration values for access and refresh token expiry
		accessTEX, err := time.ParseDuration(viper.GetString("jwt.access_token_expired"))
		if err != nil {
			log.Fatalf("Invalid access token expiry: %v", err)
		}

		refreshTEX, err := time.ParseDuration(viper.GetString("jwt.refresh_token_expired"))
		if err != nil {
			log.Fatalf("Invalid refresh token expiry: %v", err)
		}

		// Setting JWT struct values
		configInstance.JWT = &JWT{
			JWTSecret:           viper.GetString("jwt.jwt_secret"),
			AccessTokenExpired:  accessTEX,
			RefreshTokenExpired: refreshTEX,
			Issuer:              viper.GetString("jwt.issuer"),
		}
	})

	return configInstance
}
