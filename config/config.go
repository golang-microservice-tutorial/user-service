package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppEnv  string
	AppPort string
	DB      DBConfig
	JWT     JWTConfig
	Redis   RedisConfig
	Kafka   KafkaConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type JWTConfig struct {
	PrivateKeyPath string
	PublicKeyPath  string
}

type RedisConfig struct {
	Addr string
}

type KafkaConfig struct {
	Brokers []string
}

func LoadConfig() *AppConfig {
	_ = godotenv.Load()

	cfg := &AppConfig{
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "8080"),
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "app"),
		},
		JWT: JWTConfig{
			PrivateKeyPath: getEnv("JWT_PRIVATE_KEY_PATH", "key/private.pem"),
			PublicKeyPath:  getEnv("JWT_PUBLIC_KEY_PATH", "key/public.pem"),
		},
		Redis: RedisConfig{
			Addr: getEnv("REDIS_ADDR", "localhost:6379"),
		},
		Kafka: KafkaConfig{
			Brokers: strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
		},
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
