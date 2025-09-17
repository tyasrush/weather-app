package domain

import (
	"database/sql"
	"time"
)

type ForecastType string

const (
	ForecastTypeDay  ForecastType = "day"
	ForecastTypeHour ForecastType = "hour"
)

type Weather struct {
	ID                    int64        `json:"id"`
	LocationID            int64        `json:"location_id"`
	TemperatureCelcius    float64      `json:"temperature_celcius"`
	TemperatureFahrenheit float64      `json:"temperature_fahrenheit"`
	Humidity              int          `json:"humidity"`
	WindSpeed             float64      `json:"wind_speed"`
	ConditionStatus       string       `json:"condition_status"`
	ConditionIconURL      string       `json:"condition_icon_url"`
	ForecastTime          time.Time    `json:"forecast_time"`
	ForecastType          ForecastType `json:"forecast_type"`
	CreatedAt             time.Time    `json:"created_at"`
	LastModifiedAt        sql.NullTime `json:"last_modified_at"`
	DeletedAt             sql.NullTime `json:"deleted_at"`
}
