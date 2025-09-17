package config

import "os"

type Config struct {
	Port              string
	MySQLDSN          string
	RedisAddr         string
	RedisPassword     string
	WeatherAPIBaseURL string
	WeatherAPIKey     string
}

func Load() *Config {
	return &Config{
		Port:              getEnv("PORT", "8080"),
		MySQLDSN:          getEnv("MYSQL_DSN", "admin:admin@tcp(mysql:3306)/weather-db?charset=utf8mb4&parseTime=true&loc=Local"),
		RedisAddr:         getEnv("REDIS_ADDR", "redis:6379"),
		RedisPassword:     getEnv("REDIS_PASSWORD", ""),
		WeatherAPIBaseURL: getEnv("WEATHER_API_BASE_URL", "https://api.weatherapi.com/v1"),
		WeatherAPIKey:     getEnv("WEATHER_API_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
