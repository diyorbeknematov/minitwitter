package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type ServerConfig struct {
	GRPCPort int
	HTTPPort int
}

type DBConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

type Config struct {
	Server ServerConfig
	DB     DBConfig
}

func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	return &Config{
		Server: ServerConfig{
			GRPCPort: cast.ToInt(getEnv("GRPC_PORT", 50051)),
			HTTPPort: cast.ToInt(getEnv("HTTP_PORT", 8080)),
		},
		DB: DBConfig{
			Host:     cast.ToString(getEnv("DB_HOST", "localhost")),
			User:     cast.ToString(getEnv("DB_USER", "postgres")),
			Port:     cast.ToInt(getEnv("DB_PORT", 5432)),
			Name:     cast.ToString(getEnv("DB_NAME", "notification-service")),
			Password: cast.ToString(getEnv("DB_PASSWORD", "passw")),
		},
	}
}

func getEnv(key string, defaultValue interface{}) interface{} {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}

	return value
}
