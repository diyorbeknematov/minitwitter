package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	GRPCPort string

	DBHost     string
	DBName     string
	DBPort     string
	DBUser     string
	DBPassword string
}

func Load() *Config {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Error loading .env file")
	}

	return &Config{
		GRPCPort: cast.ToString(getEnv("GRPC_PORT", 50051)),

		DBHost: cast.ToString(getEnv("DB_HOST", "localhost")),
		DBName: cast.ToString(getEnv("DB_NAME", "media-service")),
		DBPort: cast.ToString(getEnv("DB_PORT", 5432)),
		DBUser: cast.ToString(getEnv("DB_USER", "postgres")),
		DBPassword: cast.ToString(getEnv("DB_PASSWORD", "pass")),
	}
}

func getEnv(key string, defaultValue interface{}) interface{} {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}

	return value
}
