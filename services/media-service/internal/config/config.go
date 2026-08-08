package config

import (
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type GRPCConfig struct {
	Port int
	Host string
}

type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type MinIOConfig struct {
	Endpoint        string
	AccessKey       string
	SecretKey       string
	Bucket          string
	UseSSL          bool
	PresignedExpiry time.Duration
}

type Config struct {
	GRPC  GRPCConfig
	DB    DBConfig
	MinIO MinIOConfig
}

func Load() *Config {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Error loading .env file")
	}

	return &Config{
		GRPC: GRPCConfig{
			Port: cast.ToInt(getEnv("_MEDIA_GRPC_PORT", 50051)),
			Host: cast.ToString(getEnv("MEDIA_GRPC_HOST", "localhost")),
		},
		DB: DBConfig{
			Host:     cast.ToString(getEnv("DB_HOST", "localhost")),
			Name:     cast.ToString(getEnv("DB_NAME", "media-service")),
			Port:     cast.ToString(getEnv("DB_PORT", 5432)),
			User:     cast.ToString(getEnv("DB_USER", "postgres")),
			Password: cast.ToString(getEnv("DB_PASSWORD", "pass")),
		},
		MinIO: MinIOConfig{
			Endpoint:        cast.ToString(getEnv("MINIO_ENDPOINT", "localhost:9000")),
			AccessKey:       cast.ToString(getEnv("MINIO_ACCESS_KEY", "minioadmin")),
			SecretKey:       cast.ToString(getEnv("MINIO_SECRET_KEY", "minioadmin")),
			Bucket:          cast.ToString(getEnv("MINIO_BUCKET", "media-service")),
			UseSSL:          cast.ToBool(getEnv("MINIO_USE_SSL", false)),
			PresignedExpiry: cast.ToDuration(getEnv("PRESIGNED_URL_EXPIRY", "72h")),
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
