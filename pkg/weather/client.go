package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"tyarus/weather-app/internal/config"
	"tyarus/weather-app/pkg/utils"
)

type WeatherAPIClientInterface interface {
	GetForecast(ctx context.Context, location string, day int) (*ForecastResponse, error)
}

type WeatherAPIClient struct {
	HTTP   *http.Client
	Config config.Config
}

func NewClient(config config.Config) WeatherAPIClientInterface {
	return &WeatherAPIClient{
		Config: config,
		HTTP: &http.Client{
			Timeout: utils.DefaultHTTPTimeout,
		},
	}
}

func (c *WeatherAPIClient) GetForecast(ctx context.Context, location string, day int) (*ForecastResponse, error) {
	if day == 0 {
		day = 14
	}

	var forecast ForecastResponse
	fetchForecastFunc := func() error {
		url := fmt.Sprintf("%s/forecast.json?key=%s&q=%s&days=%d",
			c.Config.WeatherAPIBaseURL, c.Config.WeatherAPIKey, location, day)

		resp, err := c.HTTP.Get(url)
		if err != nil {
			return fmt.Errorf("failed to request weather api: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("weather api returned status %d", resp.StatusCode)
		}

		if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}

		return nil
	}

	err := utils.RetryWithBackoff(ctx, utils.RetryWithBackoffParam{
		Func:       fetchForecastFunc,
		BaseDelay:  time.Duration(c.Config.BackoffBaseDelay),
		MaxRetries: c.Config.BackoffMaxRetries,
		MaxDelay:   time.Duration(c.Config.BackoffMaxDelay),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to fetch forecast with backoff: %w", err)
	}

	return &forecast, nil
}
