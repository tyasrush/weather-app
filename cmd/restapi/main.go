package main

import (
	"log"
	"net/http"
	"tyarus/weather-app/internal/config"
	"tyarus/weather-app/internal/handler"
	"tyarus/weather-app/internal/infra"
	"tyarus/weather-app/internal/repository"
	"tyarus/weather-app/internal/usecase"
	"tyarus/weather-app/pkg/weather"

	"github.com/gorilla/mux"
)

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

	locationUc := usecase.NewLocationUsecase(locationRepo)
	weatherUc := usecase.NewWeatherUsecase(weatherRepo, locationRepo, cache, weatherAPIClient)

	commonHandler := handler.NewCommonHandler(db, cache)
	locationHandler := handler.NewLocationHandler(locationUc)
	weatherHandler := handler.NewWeatherHandler(weatherUc)

	routes := mux.NewRouter()
	routes.HandleFunc("/health", commonHandler.HealthCheck()).Methods(http.MethodGet)
	routes.HandleFunc("/ready", commonHandler.ReadyCheck()).Methods(http.MethodGet)

	apiRoutes := routes.PathPrefix("/api/v1").Subrouter()

	apiRoutes.HandleFunc("/locations", locationHandler.GetLocationHandler()).Methods(http.MethodGet)
	apiRoutes.HandleFunc("/locations", locationHandler.CreateLocationHandler()).Methods(http.MethodPost)
	apiRoutes.HandleFunc("/weathers/sync", weatherHandler.SyncWeatherHandler()).Methods(http.MethodPost)
	apiRoutes.HandleFunc("/weathers", weatherHandler.GetWeathersHandler()).Methods(http.MethodGet)

	addr := ":" + cfg.Port
	log.Println("Server running on", addr)
	if err := http.ListenAndServe(addr, routes); err != nil {
		log.Fatal(err)
	}
}
