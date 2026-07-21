package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr              	string
	DatabaseURL       	string
	RedisURL          	string
	RSAPrivateKeyPEM  	string
	MigrationsPath    	string
	RateLimitRequests 	int
	RateLimitWindow   	time.Duration
	Env               	string
	JWTIssuer         	string
	AccessTokenTTL    	time.Duration
	RefreshTokenTTL   	time.Duration
	AllowedOrigin     	string
	TempPassword	  	string
	PrometheusUsername	string
	PrometheusPassword	string
	AccountID			string
	AccessKey			string
	SecretKey			string
	Bucket				string
	APIEndpoint			string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	private_key, err := base64.StdEncoding.DecodeString(getEnv("RSA_PRIVATE_KEY_PEM"))
	if err != nil {
		return nil, fmt.Errorf("failed to decode RSA_PRIVATE_KEY_PEM: %w", err)
	}

	cfg := &Config{
		Addr:              	":" + getEnv("PORT"),
		DatabaseURL:       	getEnv("DATABASE_URL"),
		RedisURL:          	getEnv("REDIS_URL"),
		RSAPrivateKeyPEM:  	string(private_key),
		MigrationsPath:    	getEnv("MIGRATIONS_PATH"),
		Env:               	getEnv("ENV"),
		JWTIssuer:         	getEnv("JWT_ISSUER"),
		AllowedOrigin:     	getEnv("ALLOWED_ORIGIN"),
		TempPassword: 		getEnv("TEMP_PASSWORD"),
		PrometheusUsername: getEnv("PROMETHEUS_USERNAME"),
		PrometheusPassword: getEnv("PROMETHEUS_PASSWORD"),
		AccountID:			getEnv("ACCOUNT_ID"),
		AccessKey:			getEnv("ACCESS_KEY"),
		SecretKey:			getEnv("SECRET_ACCESS_KEY"),
		Bucket: 			getEnv("BUCKET"),
		APIEndpoint:		getEnv("API_ENDPOINT"),
	}

	limitStr := getEnv("RATE_LIMIT_REQUESTS")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMIT_REQUESTS: %w", err)
	}
	cfg.RateLimitRequests = limit

	window, err := time.ParseDuration(getEnv("RATE_LIMIT_WINDOW"))
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMIT_WINDOW: %w", err)
	}
	cfg.RateLimitWindow = window

	accessTTL, err := time.ParseDuration(getEnv("ACCESS_TOKEN_TTL"))
	if err != nil {
		return nil, fmt.Errorf("invalid ACCESS_TOKEN_TTL: %w", err)
	}
	cfg.AccessTokenTTL = accessTTL

	refreshTTL, err := time.ParseDuration(getEnv("REFRESH_TOKEN_TTL")) // 30 days
	if err != nil {
		return nil, fmt.Errorf("invalid REFRESH_TOKEN_TTL: %w", err)
	}
	cfg.RefreshTokenTTL = refreshTTL

	return cfg, nil
}

func getEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("required environment variable %q is not set", key))
	}
	return v
}