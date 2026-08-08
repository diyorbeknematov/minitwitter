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

type ServiceConfig struct {
	Host string
	Port int
}

type GRPCConfig struct {
	User         ServiceConfig
	Tweet        ServiceConfig
	Media        ServiceConfig
	Notification ServiceConfig
}

type Config struct {
	Server ServerConfig
	GRPC   GRPCConfig
}

func Load() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("warning: .env file not found")
	}

	return &Config{
		Server: ServerConfig{
			Host: cast.ToString(getEnv("HTTP_HOST", "0.0.0.0")),
			Port: cast.ToInt(getEnv("HTTP_PORT", "8080")),
		},
		GRPC: GRPCConfig{
			User: ServiceConfig{
				Host: cast.ToString(getEnv("USER_GRPC_HOST", "localhost")),
				Port: cast.ToInt(getEnv("USER_GRPC_PORT", "50051")),
			},
			Tweet: ServiceConfig{
				Host: cast.ToString(getEnv("TWEET_GRPC_HOST", "localhost")),
				Port: cast.ToInt(getEnv("TWEET_GRPC_PORT", "50052")),
			},
			Media: ServiceConfig{
				Host: cast.ToString(getEnv("MEDIA_GRPC_HOST", "localhost")),
				Port: cast.ToInt(getEnv("MEDIA_GRPC_PORT", "50053")),
			},
			Notification: ServiceConfig{
				Host: cast.ToString(getEnv("NOTIFICATION_GRPC_HOST", "localhost")),
				Port: cast.ToInt(getEnv("NOTIFICATION_GRPC_PORT", "50054")),
			},
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

func (c ServiceConfig) Address() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}
