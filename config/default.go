package config

import (
	"time"

	"github.com/spf13/viper"
)

var (
	DefaultEnvironment = "dev"
)

type Config struct {
	DBUri         string `mapstructure:"MONGODB_URI"`
	DBName        string `mapstructure:"MONGODB_DB_NAME"`
	Port          string `mapstructure:"PORT"`
	RedisUri      string `mapstructure:"REDIS_URL"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	Environment   string `mapstructure:"ENVIRONMENT"`
	ProxyScanUrl  string `mapstructure:"PROXY_SCAN_URL"`
	Domain        string `mapstructure:"DOMAIN"`

	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`

	Origin string `mapstructure:"CLIENT_ORIGIN"`

	EmailFrom string `mapstructure:"EMAIL_FROM"`
	SMTPHost  string `mapstructure:"SMTP_HOST"`
	SMTPPass  string `mapstructure:"SMTP_PASS"`
	SMTPPort  int    `mapstructure:"SMTP_PORT"`
	SMTPUser  string `mapstructure:"SMTP_USER"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName(".env")

	viper.SetDefault("Environment", DefaultEnvironment)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
