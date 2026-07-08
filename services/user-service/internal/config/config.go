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

	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	AccesstokenSecret  string
	RefreshtokenSecret string
}

func ConfigLoad() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	return &Config{
		GRPCPort: cast.ToString(getEnv("GRPC_PORT", "50051")),

		DBHost: cast.ToString(getEnv("DB_HOST", "localhost")),
		DBName: cast.ToString(getEnv("DB_NAME", "user_service")),
		DBPort: cast.ToString(getEnv("DB_PORT", "5432")),
		DBUser: cast.ToString(getEnv("DB_USER", "postgres")),
		DBPassword: cast.ToString(getEnv("DB_PASSWORD", "")),

		RedisHost: cast.ToString(getEnv("REDIS_HOST", "localhost")),
		RedisPort: cast.ToString(getEnv("REDIS_PORT", "6379")),
		RedisPassword: cast.ToString(getEnv("REDIS_PASSWORD", "")),
		RedisDB: cast.ToInt(getEnv("REDIS_DB", 0)),

		AccesstokenSecret:  cast.ToString(getEnv("ACCESS_TOKEN_SECRET", "access-secret")),
		RefreshtokenSecret: cast.ToString(getEnv("REFRESH_TOKEN_SECRET", "refresh-secret")),
	}
}

func getEnv(key string, defaultValue interface{}) interface{} {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}
	return value
}
