package config

import (
	"fmt"

	"github.com/hafizhproject45/Golang-Boilerplate.git/internal/utils"

	"github.com/spf13/viper"
)

var (
	IsProd              bool
	AppHost             string
	Version             string
	LogLevel            string
	AppPort             int
	DBHost              string
	DBUser              string
	DBPassword          string
	DBName              string
	DBPort              int
	JWTSecret           string
	JWTAccessExp        int
	JWTRefreshExp       int
	JWTResetPasswordExp int
	JWTVerifyEmailExp   int
	PostgresDSN         string
	RedisURL            string
	Issuer              string
	SMTPHost            string
	SMTPPort            int
	SMTPUsername        string
	SMTPPassword        string
	EmailFrom           string
	GoogleClientID      string
	GoogleClientSecret  string
	RedirectURL         string
)

func init() {
	loadConfig()

	// server configuration
	IsProd = viper.GetString("APP_ENV") == "prod"
	// AppHost = viper.GetString("APP_HOST")
	// AppPort = viper.GetInt("APP_PORT")
	AppHost = viper.GetString("APP_HOST")
	if AppHost == "" {
		AppHost = "0.0.0.0"
	}
	AppPort = viper.GetInt("APP_PORT")
	if AppPort == 0 {
		AppPort = 8080
	}
	Version = viper.GetString("VERSION")
	LogLevel = viper.GetString("LOG_LEVEL")

	// database configuration
	DBHost = viper.GetString("DB_HOST")
	DBUser = viper.GetString("DB_USER")
	DBPassword = viper.GetString("DB_PASSWORD")
	DBName = viper.GetString("DB_NAME")
	DBPort = viper.GetInt("DB_PORT")
	PostgresDSN = viper.GetString("POSTGRES_DSN")
	if PostgresDSN == "" {
		PostgresDSN = fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=disable",
			DBUser, DBPassword, DBHost, DBPort, DBName,
		)
	}

	// jwt configuration
	JWTSecret = viper.GetString("JWT_SECRET")
	JWTAccessExp = viper.GetInt("JWT_ACCESS_EXP_MINUTES")
	JWTRefreshExp = viper.GetInt("JWT_REFRESH_EXP_DAYS")
	JWTResetPasswordExp = viper.GetInt("JWT_RESET_PASSWORD_EXP_MINUTES")
	JWTVerifyEmailExp = viper.GetInt("JWT_VERIFY_EMAIL_EXP_MINUTES")

	// Redis / OIDC
	RedisURL = viper.GetString("REDIS_URL")
	if RedisURL == "" {
		RedisURL = "redis://redis:6379/0"
	}
	Issuer = viper.GetString("ISSUER")
	if Issuer == "" {
		// fallback ke SSO_ISSUER jika kamu sudah pakai itu sebelumnya
		Issuer = viper.GetString("SSO_ISSUER")
	}
	// SMTP configuration
	SMTPHost = viper.GetString("SMTP_HOST")
	SMTPPort = viper.GetInt("SMTP_PORT")
	SMTPUsername = viper.GetString("SMTP_USERNAME")
	SMTPPassword = viper.GetString("SMTP_PASSWORD")
	EmailFrom = viper.GetString("EMAIL_FROM")

	// oauth2 configuration
	GoogleClientID = viper.GetString("GOOGLE_CLIENT_ID")
	GoogleClientSecret = viper.GetString("GOOGLE_CLIENT_SECRET")
	RedirectURL = viper.GetString("REDIRECT_URL")
}

func loadConfig() {
	viper.AutomaticEnv()

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err == nil {
		utils.Log.Info("Config file loaded from .env")
	} else {
		utils.Log.Warn("No .env file found, using environment variables only")
	}
}
