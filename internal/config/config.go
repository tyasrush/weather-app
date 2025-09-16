package config

import "os"

type Config struct {
	Port          string
	MySQLDSN      string
	RedisAddr     string
	RedisPassword string
}

func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "8080"),
		MySQLDSN:      getEnv("MYSQL_DSN", "admin:admin@tcp(mysql:3306)/weather-db"),
		RedisAddr:     getEnv("REDIS_ADDR", "redis:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
