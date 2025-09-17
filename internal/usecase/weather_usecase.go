package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
	"tyarus/weather-app/internal/domain"
	"tyarus/weather-app/internal/dto"
	"tyarus/weather-app/internal/repository"
	"tyarus/weather-app/pkg/response"
	"tyarus/weather-app/pkg/utils"
	"tyarus/weather-app/pkg/weather"

	"github.com/redis/go-redis/v9"
)

type WeatherUsecaseInterface interface {
	SyncWeatherUsecase(ctx context.Context, req dto.PostWeatherSyncUsecaseRequest) error
	GetWeathersUsecase(ctx context.Context, req dto.GetWeathersParam) (response.Response[dto.GetWeatherResponse], error)
}

type weatherUsecase struct {
	weatherRepo      repository.WeatherRepositoryInterface
	locationRepo     repository.LocationRepositoryInterface
	cache            *redis.Client
	weatherAPIClient weather.WeatherAPIClientInterface
}

func NewWeatherUsecase(
	weatherRepo repository.WeatherRepositoryInterface,
	locationRepo repository.LocationRepositoryInterface,
	cache *redis.Client,
	weatherAPIClient weather.WeatherAPIClientInterface,
) WeatherUsecaseInterface {
	return &weatherUsecase{
		weatherRepo:      weatherRepo,
		locationRepo:     locationRepo,
		cache:            cache,
		weatherAPIClient: weatherAPIClient,
	}
}

func (u *weatherUsecase) SyncWeatherUsecase(ctx context.Context, req dto.PostWeatherSyncUsecaseRequest) error {
	param := repository.GetLocationsParam{}
	if req.Limit == 0 {
		req.Limit = 10
	}
	param.Limit = req.Limit

	if req.LocationID > 0 {
		param.ID = req.LocationID
	}

	locations, err := u.locationRepo.GetLocations(ctx, param)
	if err != nil {
		return fmt.Errorf("failed to get locations: %w", err)
	}

	if req.ForecastDayTotal == 0 {
		req.ForecastDayTotal = 14 // max day from weather api
	}

	for _, location := range locations {
		forecast, err := u.weatherAPIClient.GetForecast(location.Name, req.ForecastDayTotal)
		if err != nil {
			log.Printf("Failed to get forecast for location %s: %v", location.Name, err)
			continue
		}

		var weathers []domain.Weather
		for _, day := range forecast.Forecast.Forecastday {
			fmt.Println("date: ", day.Date)
			forecastTime, err := time.Parse(utils.DateFormat, day.Date)
			if err != nil {
				log.Printf("Failed to parse forecast date %s: %v", day.Date, err)
				continue
			}

			forecastWeather := domain.Weather{
				LocationID:            location.ID,
				TemperatureCelcius:    day.Day.AvgtempC,
				TemperatureFahrenheit: day.Day.AvgtempF,
				Humidity:              int(day.Day.AvgHumidity),
				WindSpeed:             day.Day.MaxWindKPH,
				ConditionStatus:       day.Day.Condition.Text,
				ConditionIconURL:      day.Day.Condition.Icon,
				ForecastTime:          forecastTime,
				ForecastType:          domain.ForecastTypeDay,
			}

			weathers = append(weathers, forecastWeather)

			for _, item := range day.Hours {
				forecastHourTime, err := time.Parse(utils.DateFormatWithHour, item.ForecastTime)
				if err != nil {
					log.Printf("Failed to parse forecast hour time %s: %v", item.ForecastTime, err)
					continue
				}

				forecastWeather := domain.Weather{
					LocationID:            location.ID,
					TemperatureCelcius:    item.TempC,
					TemperatureFahrenheit: item.TempF,
					Humidity:              int(item.Humidity),
					WindSpeed:             item.WindKph,
					ConditionStatus:       item.Condition.Text,
					ConditionIconURL:      item.Condition.Icon,
					ForecastTime:          forecastHourTime,
					ForecastType:          domain.ForecastTypeHour,
				}

				weathers = append(weathers, forecastWeather)
			}
		}

		_, err = u.weatherRepo.BulkUpsertWeather(ctx, weathers)
		if err != nil {
			log.Printf("Failed to bulk upsert weather data %s: %v", location.Name, err)
			continue
		}

		log.Printf("sync weather data success: %s", location.Name)
	}

	return nil
}

func (u *weatherUsecase) GetWeathersUsecase(ctx context.Context, param dto.GetWeathersParam) (response.Response[dto.GetWeatherResponse], error) {
	resp := response.Response[dto.GetWeatherResponse]{
		Status:  "success",
		Message: "get weather data success",
	}
	locations, err := u.locationRepo.GetLocations(ctx, repository.GetLocationsParam{
		ID:    int(param.LocationID),
		Limit: 1,
	})
	if err != nil {
		return resp, fmt.Errorf("failed to get locations: %w", err)
	}

	if len(locations) == 0 {
		return resp, errors.New("location not found, please check your parameter")
	}

	cacheKey := fmt.Sprintf(utils.WeatherLocationKey, param.LocationID)
	cachedData, err := u.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		var weatherResponse dto.GetWeatherResponse
		if err := json.Unmarshal([]byte(cachedData), &weatherResponse); err == nil {
			resp.Data = weatherResponse
			return resp, nil
		}
	}

	offset := (param.CurrentPage - 1) * param.PageSize
	repoParam := repository.GetWeathersParam{
		LocationID: int64(param.LocationID),
		Limit:      param.PageSize,
		Offset:     offset,
		OrderBy:    "forecast_time DESC",
	}

	weathers, err := u.weatherRepo.GetWeathers(ctx, repoParam)
	if err != nil {
		return resp, fmt.Errorf("failed to get weathers: %w", err)
	}

	if len(weathers) == 0 {
		return resp, nil
	}

	location := locations[0]
	locationResponse := dto.GetLocationHandlerResponseItem{
		ID:        int64(param.LocationID),
		Name:      location.Name,
		Region:    location.Region,
		Country:   location.Country,
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
		CreatedAt: time.Now(),
	}

	forecast := []dto.GetWeatherResponseItem{}
	currentTimeWeather := dto.GetWeatherResponseItem{}
	var currentTimeDuration time.Duration = 0
	for _, item := range weathers[1:] {
		itemResponse := dto.GetWeatherResponseItem{
			ForecastTime:          item.ForecastTime,
			ForecastType:          string(item.ForecastType),
			TemperatureCelcius:    item.TemperatureCelcius,
			TemperatureFahrenheit: item.TemperatureFahrenheit,
			Humidity:              item.Humidity,
			WindSpeed:             item.WindSpeed,
			Condition: dto.WeatherConditionResponse{
				Status:  item.ConditionStatus,
				IconURL: item.ConditionIconURL,
			},
			CreatedAt:      item.CreatedAt,
			LastModifiedAt: item.LastModifiedAt.Time,
		}

		if currentTimeDuration == 0 ||
			(time.Now().Day() == item.ForecastTime.Day() &&
				time.Now().Sub(item.ForecastTime) < currentTimeDuration) {
			currentTimeWeather = itemResponse
			currentTimeDuration = time.Now().Sub(item.ForecastTime)
		}

		forecast = append(forecast, itemResponse)
	}

	weatherResponse := dto.GetWeatherResponse{
		Location:    locationResponse,
		CurrentTime: currentTimeWeather,
		Forecast:    forecast,
	}

	cacheData, err := json.Marshal(weatherResponse)
	if err == nil {
		err = u.cache.Set(ctx, cacheKey, cacheData, 10*time.Minute).Err()
		if err != nil {
			log.Printf("Failed to set cache: %v", err)
		}
	}

	resp.Data = weatherResponse
	return resp, nil
}
