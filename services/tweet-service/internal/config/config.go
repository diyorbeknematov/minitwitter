package config

import (
	"log"
	"net"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type GRPCConfig struct {
	Host string
	Port int
}

type DBConfig struct {
	Host     string
	Name     string
	Port     string
	User     string
	Password string
}

type Config struct {
	GRPC GRPCConfig
	DB   DBConfig
}

func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	return &Config{
		GRPC: GRPCConfig{
			Host: cast.ToString(getEnv("TWEET_GRPC_HOST", "localhost")),
			Port: cast.ToInt(getEnv("MEDIA_GRPC_PORT", 50053)),
		},

		DB: DBConfig{
			Host:     cast.ToString(getEnv("DB_HOST", "localhost")),
			Name:     cast.ToString(getEnv("DB_NAME", "user_service")),
			Port:     cast.ToString(getEnv("DB_PORT", "5432")),
			User:     cast.ToString(getEnv("DB_USER", "postgres")),
			Password: cast.ToString(getEnv("DB_PASSWORD", "")),
		},
	}
}

func (c GRPCConfig) Address() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

func getEnv(key string, defaultValue interface{}) interface{} {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}
	return value
}
