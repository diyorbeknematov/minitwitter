package config

import (
	"log"
	"net"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type ServerConfig struct {
	Host string
	Port int
}

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

type Config struct {
	Server ServerConfig
	GRPC   GRPCConfig
	DB     DBConfig
}

func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	return &Config{
		Server: ServerConfig{
			Host: cast.ToString(getEnv("HTTP_HOST", 8080)),
			Port: cast.ToInt(getEnv("HTTP_PORT", 50051)),
		},
		GRPC: GRPCConfig{
			Host: cast.ToString(getEnv("NOTIFICATION_GRPC_HOST", "localhost")),
			Port: cast.ToInt(getEnv("NOTIFICATION_GRPC_PORT", 50054)),
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

func (c GRPCConfig) Address() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

func (c ServerConfig) Address() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

func getEnv(key string, defaultValue interface{}) interface{} {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}

	return value
}
