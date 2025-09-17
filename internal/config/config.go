package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port              string
	MySQLDSN          string
	RedisAddr         string
	RedisPassword     string
	WeatherAPIBaseURL string
	WeatherAPIKey     string
	BackoffMaxRetries int
	BackoffBaseDelay  int
	BackoffMaxDelay   int
	WorkerPeriod      int
	WorkerLimit       int
}

func Load() *Config {
	return &Config{
		Port:              getEnv("PORT", "8080"),
		MySQLDSN:          getEnv("MYSQL_DSN", "admin:admin@tcp(localhost:3306)/weather-db?charset=utf8mb4&parseTime=true&loc=Local"),
		RedisAddr:         getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:     getEnv("REDIS_PASSWORD", ""),
		WeatherAPIBaseURL: getEnv("WEATHER_API_BASE_URL", "https://api.weatherapi.com/v1"),
		WeatherAPIKey:     getEnv("WEATHER_API_KEY", ""),
		BackoffMaxRetries: getEnvInt("BACKOFF_MAX_RETRIES", "3"),
		BackoffBaseDelay:  getEnvInt("BACKOFF_BASE_DELAY", "200000000"),
		BackoffMaxDelay:   getEnvInt("BACKOFF_MAX_DELAY", "5000000000"),
		WorkerPeriod:      getEnvInt("WORKER_PERIOD", "900000000000"),
		WorkerLimit:       getEnvInt("WORKER_LIMIT", "10"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func getEnvInt(key, fallback string) int {
	result := fallback
	if val, ok := os.LookupEnv(key); ok {
		result = val
	}

	resultInt, _ := strconv.Atoi(result)
	return resultInt
}
