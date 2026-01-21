package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort     string
	GrpcPort    string
	AppEnv      string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DatabaseDSN string

	JWTSecret     string
	JWTExpiration int // in seconds
}

func LoadConfig(env string) *Config {

	envFile := ".env"
	if env != "" {
		envFile = ".env." + env
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Println("No .env file found, using system env", err)
	}

	jwtExp := getEnvAsInt("JWT_EXPIRATION", 3600)

	cfg := &Config{
		AppPort:       getEnv("APP_PORT", "8000"),
		GrpcPort:      getEnv("GRPC_PORT", "50052"),
		AppEnv:        getEnv("APP_ENV", "development"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", ""),
		DBName:        getEnv("DB_NAME", "test"),
		JWTSecret:     getEnv("JWT_SECRET", "changeme"),
		JWTExpiration: jwtExp,
	}

	cfg.DatabaseDSN = fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	return cfg
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	return fallback
}
