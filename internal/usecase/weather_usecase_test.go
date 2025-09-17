package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"tyarus/weather-app/internal/domain"
	"tyarus/weather-app/internal/dto"
	"tyarus/weather-app/internal/repository"
	"tyarus/weather-app/mocks"
)

func TestGetWeathersUsecase(t *testing.T) {
	t.Run("WHEN error occurred on get locations, THEN should return error accordingly", func(t *testing.T) {
		mockWeatherRepo := mocks.NewWeatherRepositoryInterface(t)
		mockLocationRepo := mocks.NewLocationRepositoryInterface(t)
		mockCache := mocks.NewCacheInterface(t)

		usecase := NewWeatherUsecase(mockWeatherRepo, mockLocationRepo, mockCache, nil)
		ctx := context.Background()
		req := dto.GetWeathersParam{
			LocationID:  1,
			PageSize:    10,
			CurrentPage: 1,
		}

		expectedError := errors.New("database error")
		mockLocationRepo.On("GetLocations", ctx, repository.GetLocationsParam{
			ID:    1,
			Limit: 1,
		}).Return(nil, expectedError)

		_, err := usecase.GetWeathersUsecase(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get locations")
	})

	t.Run("WHEN location not found, THEN should return error accordingly", func(t *testing.T) {
		mockWeatherRepo := mocks.NewWeatherRepositoryInterface(t)
		mockLocationRepo := mocks.NewLocationRepositoryInterface(t)
		mockCache := mocks.NewCacheInterface(t)

		usecase := NewWeatherUsecase(mockWeatherRepo, mockLocationRepo, mockCache, nil)
		ctx := context.Background()
		req := dto.GetWeathersParam{
			LocationID:  1,
			PageSize:    10,
			CurrentPage: 1,
		}

		locations := []domain.Location{}
		mockLocationRepo.On("GetLocations", ctx, repository.GetLocationsParam{
			ID:    1,
			Limit: 1,
		}).Return(locations, nil)

		_, err := usecase.GetWeathersUsecase(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "location not found")
	})

	t.Run("WHEN data available on cache, THEN should return data accordingly", func(t *testing.T) {
		mockWeatherRepo := mocks.NewWeatherRepositoryInterface(t)
		mockLocationRepo := mocks.NewLocationRepositoryInterface(t)
		mockCache := mocks.NewCacheInterface(t)

		usecase := NewWeatherUsecase(mockWeatherRepo, mockLocationRepo, mockCache, nil)
		ctx := context.Background()
		req := dto.GetWeathersParam{
			LocationID:  1,
			PageSize:    10,
			CurrentPage: 1,
		}

		locations := []domain.Location{
			{
				ID:        1,
				Name:      "Test Location",
				Region:    "Test Region",
				Country:   "Test Country",
				Latitude:  1.0,
				Longitude: 1.0,
			},
		}
		mockLocationRepo.On("GetLocations", ctx, repository.GetLocationsParam{
			ID:    1,
			Limit: 1,
		}).Return(locations, nil)

		cacheKey := "weather:location:1"
		expectedResponse := dto.GetWeatherResponse{
			Location: dto.GetLocationHandlerResponseItem{
				ID: 1,
			},
		}
		cacheData, _ := json.Marshal(expectedResponse)
		mockCache.On("Get", ctx, cacheKey).Return(string(cacheData), nil)

		result, err := usecase.GetWeathersUsecase(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, result.Data)
	})

	t.Run("WHEN error occurred on get weathers, THEN should return error accordingly", func(t *testing.T) {
		mockWeatherRepo := mocks.NewWeatherRepositoryInterface(t)
		mockLocationRepo := mocks.NewLocationRepositoryInterface(t)
		mockCache := mocks.NewCacheInterface(t)

		usecase := NewWeatherUsecase(mockWeatherRepo, mockLocationRepo, mockCache, nil)
		ctx := context.Background()
		req := dto.GetWeathersParam{
			LocationID:  1,
			PageSize:    10,
			CurrentPage: 1,
		}

		locations := []domain.Location{
			{
				ID:        1,
				Name:      "Test Location",
				Region:    "Test Region",
				Country:   "Test Country",
				Latitude:  1.0,
				Longitude: 1.0,
			},
		}
		mockLocationRepo.On("GetLocations", ctx, repository.GetLocationsParam{
			ID:    1,
			Limit: 1,
		}).Return(locations, nil)

		cacheKey := "weather:location:1"
		mockCache.On("Get", ctx, cacheKey).Return("", errors.New("cache miss"))

		expectedError := errors.New("database error")
		mockWeatherRepo.On("GetWeathers", ctx, repository.GetWeathersParam{
			LocationID: 1,
			Limit:      10,
			Offset:     0,
			OrderBy:    "forecast_time DESC",
		}).Return(nil, expectedError)

		_, err := usecase.GetWeathersUsecase(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get weathers")
	})

	t.Run("WHEN no weathers found, THEN should return empty forecast", func(t *testing.T) {
		mockWeatherRepo := mocks.NewWeatherRepositoryInterface(t)
		mockLocationRepo := mocks.NewLocationRepositoryInterface(t)
		mockCache := mocks.NewCacheInterface(t)

		usecase := NewWeatherUsecase(mockWeatherRepo, mockLocationRepo, mockCache, nil)
		ctx := context.Background()
		req := dto.GetWeathersParam{
			LocationID:  1,
			PageSize:    10,
			CurrentPage: 1,
		}

		locations := []domain.Location{
			{
				ID:        1,
				Name:      "Test Location",
				Region:    "Test Region",
				Country:   "Test Country",
				Latitude:  1.0,
				Longitude: 1.0,
			},
		}
		mockLocationRepo.On("GetLocations", ctx, repository.GetLocationsParam{
			ID:    1,
			Limit: 1,
		}).Return(locations, nil)

		cacheKey := "weather:location:1"
		mockCache.On("Get", ctx, cacheKey).Return("", errors.New("cache miss"))

		weathers := []domain.Weather{}
		mockWeatherRepo.On("GetWeathers", ctx, repository.GetWeathersParam{
			LocationID: 1,
			Limit:      10,
			Offset:     0,
			OrderBy:    "forecast_time DESC",
		}).Return(weathers, nil)

		result, err := usecase.GetWeathersUsecase(ctx, req)

		assert.NoError(t, err)
		assert.Empty(t, result.Data.Forecast)
	})

	t.Run("WHEN weathers found, THEN should return result accordingly", func(t *testing.T) {
		mockWeatherRepo := mocks.NewWeatherRepositoryInterface(t)
		mockLocationRepo := mocks.NewLocationRepositoryInterface(t)
		mockCache := mocks.NewCacheInterface(t)

		usecase := NewWeatherUsecase(mockWeatherRepo, mockLocationRepo, mockCache, nil)
		ctx := context.Background()
		req := dto.GetWeathersParam{
			LocationID:  1,
			PageSize:    10,
			CurrentPage: 1,
		}

		locations := []domain.Location{
			{
				ID:        1,
				Name:      "Test Location",
				Region:    "Test Region",
				Country:   "Test Country",
				Latitude:  1.0,
				Longitude: 1.0,
			},
		}
		mockLocationRepo.On("GetLocations", ctx, repository.GetLocationsParam{
			ID:    1,
			Limit: 1,
		}).Return(locations, nil)

		cacheKey := "weather:location:1"
		mockCache.On("Get", ctx, cacheKey).Return("", errors.New("cache miss"))

		now := time.Now()
		weathers := []domain.Weather{
			{
				ID:             1,
				LocationID:     1,
				ForecastTime:   now,
				ForecastType:   domain.ForecastTypeHour,
				CreatedAt:      now,
				LastModifiedAt: sql.NullTime{Valid: true, Time: now},
			},
		}
		mockWeatherRepo.On("GetWeathers", ctx, repository.GetWeathersParam{
			LocationID: 1,
			Limit:      10,
			Offset:     0,
			OrderBy:    "forecast_time DESC",
		}).Return(weathers, nil)

		mockCache.On("Set", ctx, cacheKey, mock.Anything, 10*time.Minute).Return(nil)

		result, err := usecase.GetWeathersUsecase(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), result.Data.Location.ID)
		assert.Empty(t, result.Data.Forecast)
	})

}
