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
	Port     int
	Name     string
	User     string
	Password string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type Config struct {
	GRPC  GRPCConfig
	DB    DBConfig
	Redis RedisConfig

	AccesstokenSecret  string
	RefreshtokenSecret string
}

func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	return &Config{
		GRPC: GRPCConfig{
			Host: cast.ToString(getEnv("USER_GRPC_HOST", "localhost")),
			Port: cast.ToInt(getEnv("USER_GRPC_PORT", 50051)),
		},
		DB: DBConfig{
			Host:     cast.ToString(getEnv("DB_HOST", "localhost")),
			Port:     cast.ToInt(getEnv("DB_PORT", 5432)),
			User:     cast.ToString(getEnv("DB_USER", "postgres")),
			Password: cast.ToString(getEnv("DB_PASSWORD", "pass")),
			Name:     cast.ToString(getEnv("DB_NAME", "user_service")),
		},

		Redis: RedisConfig{
			Host:     cast.ToString(getEnv("REDIS_HOST", "localhost")),
			Port:     cast.ToString(getEnv("REDIS_PORT", "6379")),
			Password: cast.ToString(getEnv("REDIS_PASSWORD", "")),
			DB:       cast.ToInt(getEnv("REDIS_DB", 0)),
		},

		AccesstokenSecret:  cast.ToString(getEnv("ACCESS_TOKEN_SECRET", "access-secret")),
		RefreshtokenSecret: cast.ToString(getEnv("REFRESH_TOKEN_SECRET", "refresh-secret")),
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
