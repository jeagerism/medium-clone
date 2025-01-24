package config

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
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

		accessTEX, err := time.ParseDuration(viper.GetString("jwt.access_token_expired"))
		if err != nil {
			log.Fatalf("Invalid access token expiry: %v", err)
		}

		refreshTEX, err := time.ParseDuration(viper.GetString("jwt.refresh_token_expired"))
		if err != nil {
			log.Fatalf("Invalid refresh token expiry: %v", err)
		}

		configInstance = &Config{
			server: &server{
				port: viper.GetInt("server.port"),
			},
			db: &db{
				host:     viper.GetString("db.host"),
				port:     viper.GetInt("db.port"),
				user:     viper.GetString("db.user"),
				password: viper.GetString("db.password"),
				dbname:   viper.GetString("db.db_name"),
				sslMode:  viper.GetString("db.sslmode"),
				timezone: viper.GetString("db.timezone"),
			},
			jwt: &jwt{
				SecretKey:           viper.GetString("jwt.jwt_secret"),
				AccessTokenExpired:  accessTEX,
				RefreshTokenExpired: refreshTEX,
				Issuer:              viper.GetString("jwt.issuer"),
			},
		}
	})

	if configInstance == nil {
		log.Fatalf("Failed to initialize configuration")
	}

	return configInstance
}
