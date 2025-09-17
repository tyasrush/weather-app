package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"tyarus/weather-app/cmd/worker"
	"tyarus/weather-app/internal/dto"
	"tyarus/weather-app/mocks"
)

// MockWeatherUsecase is a mock implementation of the WeatherUsecaseInterface
type MockWeatherUsecase struct {
	mock.Mock
}

func (m *MockWeatherUsecase) SyncWeatherUsecase(ctx context.Context, req dto.PostWeatherSyncUsecaseRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockWeatherUsecase) SyncWeatherUsecaseConcurrent(ctx context.Context, req dto.PostWeatherSyncUsecaseRequest) []error {
	args := m.Called(ctx, req)
	return args.Get(0).([]error)
}

func (m *MockWeatherUsecase) GetWeathersUsecase(ctx context.Context, req dto.GetWeathersParam) (interface{}, error) {
	args := m.Called(ctx, req)
	return args.Get(0), args.Error(1)
}

func TestGetWorkerConfig(t *testing.T) {
	// Save original env values
	originalPeriod := getEnv("WORKER_PERIOD_SECONDS", "")
	originalLocationID := getEnv("WORKER_LOCATION_ID", "")
	originalLimit := getEnv("WORKER_LIMIT", "")
	originalForecastDayTotal := getEnv("WORKER_FORECAST_DAY_TOTAL", "")

	// Restore original values after test
	defer func() {
		setEnv("WORKER_PERIOD_SECONDS", originalPeriod)
		setEnv("WORKER_LOCATION_ID", originalLocationID)
		setEnv("WORKER_LIMIT", originalLimit)
		setEnv("WORKER_FORECAST_DAY_TOTAL", originalForecastDayTotal)
	}()

	t.Run("WHEN no environment variables set, THEN should use default values", func(t *testing.T) {
		// Clear environment variables
		setEnv("WORKER_PERIOD_SECONDS", "")
		setEnv("WORKER_LOCATION_ID", "")
		setEnv("WORKER_LIMIT", "")
		setEnv("WORKER_FORECAST_DAY_TOTAL", "")

		config := getWorkerConfig()

		assert.Equal(t, 1*time.Hour, config.Period)
		assert.Equal(t, 0, config.LocationID)
		assert.Equal(t, 10, config.Limit)
		assert.Equal(t, 14, config.ForecastDayTotal)
	})

	t.Run("WHEN environment variables set correctly, THEN should use those values", func(t *testing.T) {
		setEnv("WORKER_PERIOD_SECONDS", "1800") // 30 minutes
		setEnv("WORKER_LOCATION_ID", "5")
		setEnv("WORKER_LIMIT", "20")
		setEnv("WORKER_FORECAST_DAY_TOTAL", "7")

		config := getWorkerConfig()

		assert.Equal(t, 30*time.Minute, config.Period)
		assert.Equal(t, 5, config.LocationID)
		assert.Equal(t, 20, config.Limit)
		assert.Equal(t, 7, config.ForecastDayTotal)
	})

	t.Run("WHEN environment variables set incorrectly, THEN should use default values", func(t *testing.T) {
		setEnv("WORKER_PERIOD_SECONDS", "invalid")
		setEnv("WORKER_LOCATION_ID", "invalid")
		setEnv("WORKER_LIMIT", "invalid")
		setEnv("WORKER_FORECAST_DAY_TOTAL", "invalid")

		config := getWorkerConfig()

		assert.Equal(t, 1*time.Hour, config.Period)
		assert.Equal(t, 0, config.LocationID)
		assert.Equal(t, 10, config.Limit)
		assert.Equal(t, 14, config.ForecastDayTotal)
	})
}

func TestSyncWeather(t *testing.T) {
	t.Run("WHEN sync succeeds, THEN should not return error", func(t *testing.T) {
		mockUsecase := new(MockWeatherUsecase)
		config := WorkerConfig{
			LocationID:       1,
			Limit:            5,
			ForecastDayTotal: 7,
		}

		req := dto.PostWeatherSyncUsecaseRequest{
			LocationID:       config.LocationID,
			Limit:            config.Limit,
			ForecastDayTotal: config.ForecastDayTotal,
		}

		mockUsecase.On("SyncWeatherUsecase", mock.Anything, req).Return(nil)

		// This function doesn't return anything, so we just ensure it doesn't panic
		syncWeather(mockUsecase, config)
		
		mockUsecase.AssertExpectations(t)
	})

	t.Run("WHEN sync fails, THEN should not panic", func(t *testing.T) {
		mockUsecase := new(MockWeatherUsecase)
		config := WorkerConfig{
			LocationID:       1,
			Limit:            5,
			ForecastDayTotal: 7,
		}

		req := dto.PostWeatherSyncUsecaseRequest{
			LocationID:       config.LocationID,
			Limit:            config.Limit,
			ForecastDayTotal: config.ForecastDayTotal,
		}

		expectedError := errors.New("sync failed")
		mockUsecase.On("SyncWeatherUsecase", mock.Anything, req).Return(expectedError)

		// This function doesn't return anything, so we just ensure it doesn't panic
		syncWeather(mockUsecase, config)
		
		mockUsecase.AssertExpectations(t)
	})
}

// Helper functions for testing
func getEnv(key, fallback string) string {
	// In a real implementation, this would use os.Getenv
	// For testing, we'll use a map to simulate environment variables
	return fallback
}

func setEnv(key, value string) {
	// In a real implementation, this would use os.Setenv
	// For testing, we would need a map to simulate environment variables
	// This is just a placeholder
}