package main

import (
	"context"
	"log"
	"time"

	"tyarus/weather-app/internal/config"
	"tyarus/weather-app/internal/dto"
	"tyarus/weather-app/internal/infra"
	"tyarus/weather-app/internal/repository"
	"tyarus/weather-app/internal/usecase"
	"tyarus/weather-app/pkg/weather"
)

func runWorker(weatherUsecase usecase.WeatherUsecaseInterface, config config.Config) {
	log.Printf("start weather sync worker with period: %v\n", config.WorkerPeriod)

	ticker := time.NewTicker(time.Duration(config.WorkerPeriod))
	defer ticker.Stop()

	syncWeather(weatherUsecase, config)

	for {
		select {
		case <-ticker.C:
			syncWeather(weatherUsecase, config)
		case <-context.Background().Done():
			log.Println("worker stopped")
			return
		}
	}
}

func syncWeather(weatherUsecase usecase.WeatherUsecaseInterface, config config.Config) {
	log.Println("start weather sync")

	ctx := context.Background()
	req := dto.PostWeatherSyncUsecaseRequest{
		Limit: config.WorkerLimit,
	}

	err := weatherUsecase.SyncWeatherUsecase(ctx, req)
	if err != nil {
		log.Printf("failed to sync weather: %v", err)
	} else {
		log.Println("sync weather completed")
	}
}

func main() {
	cfg := config.Load()

	db, err := infra.InitDatabase(cfg.MySQLDSN)
	if err != nil {
		log.Fatalf("failed to connect MySQL: %v", err)
	}
	defer db.Close()

	cache := infra.InitCache(cfg.RedisAddr, cfg.RedisPassword)
	defer cache.Close()

	weatherAPIClient := weather.NewClient(*cfg)

	locationRepo := repository.NewLocationRepository(db)
	weatherRepo := repository.NewWeatherRepository(db)

	weatherUsecase := usecase.NewWeatherUsecase(weatherRepo, locationRepo, cache, weatherAPIClient)

	runWorker(weatherUsecase, *cfg)
}
