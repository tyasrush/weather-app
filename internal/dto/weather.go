package dto

import (
	"time"
)

type GetWeatherResponseItem struct {
	ForecastTime          time.Time                `json:"forecastTime"`
	ForecastType          string                   `json:"forecastType"`
	TemperatureCelcius    float64                  `json:"temperatureCelcius"`
	TemperatureFahrenheit float64                  `json:"temperatureFahrenheit"`
	Temperature           float64                  `json:"temperature"`
	Humidity              int                      `json:"humidity"`
	WindSpeed             float64                  `json:"windSpeed"`
	Condition             WeatherConditionResponse `json:"condition"`
	CreatedAt             time.Time                `json:"createdAt"`
	LastModifiedAt        time.Time                `json:"lastModifiedAt"`
}

type WeatherConditionResponse struct {
	Status  string `json:"status"`
	IconURL string `json:"iconURL"`
}

type GetWeatherResponse struct {
	Location    GetLocationHandlerResponseItem `json:"location,omitempty"`
	CurrentTime GetWeatherResponseItem         `json:"currentTime,omitempty"`
	Forecast    []GetWeatherResponseItem       `json:"forecast,omitempty"`
}

type GetWeathersParam struct {
	LocationID  int
	PageSize    int
	CurrentPage int
}

type PostWeatherSyncUsecaseRequest struct {
	LocationID       int `json:"locationID"`
	Limit            int `json:"limit"`
	ForecastDayTotal int `json:"forecastDayTotal"`
}
